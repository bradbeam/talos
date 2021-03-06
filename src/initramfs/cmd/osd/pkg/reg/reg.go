//+build linux

package reg

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/autonomy/talos/src/initramfs/cmd/init/pkg/constants"
	"github.com/autonomy/talos/src/initramfs/cmd/init/pkg/system/runner"
	containerdrunner "github.com/autonomy/talos/src/initramfs/cmd/init/pkg/system/runner/containerd"
	servicelog "github.com/autonomy/talos/src/initramfs/cmd/init/pkg/system/runner/process/log"
	"github.com/autonomy/talos/src/initramfs/cmd/osd/proto"
	"github.com/autonomy/talos/src/initramfs/pkg/chunker"
	filechunker "github.com/autonomy/talos/src/initramfs/pkg/chunker/file"
	streamchunker "github.com/autonomy/talos/src/initramfs/pkg/chunker/stream"
	"github.com/autonomy/talos/src/initramfs/pkg/userdata"
	"github.com/autonomy/talos/src/initramfs/pkg/version"
	"github.com/containerd/cgroups"
	"github.com/containerd/containerd"
	tasks "github.com/containerd/containerd/api/services/tasks/v1"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/typeurl"
	dockerclient "github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/golang/protobuf/ptypes/empty"
	crioclient "github.com/kubernetes-incubator/cri-o/client"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
)

// Registrator is the concrete type that implements the factory.Registrator and
// proto.OSDServer interfaces.
type Registrator struct {
	Data *userdata.UserData
}

// Register implements the factory.Registrator interface.
func (r *Registrator) Register(s *grpc.Server) {
	proto.RegisterOSDServer(s, r)
}

// Kubeconfig implements the proto.OSDServer interface. The admin kubeconfig is
// generated by kubeadm and placed at /etc/kubernetes/admin.conf. This method
// returns the contents of the generated admin.conf in the response.
func (r *Registrator) Kubeconfig(ctx context.Context, in *empty.Empty) (data *proto.Data, err error) {
	fileBytes, err := ioutil.ReadFile("/etc/kubernetes/admin.conf")
	if err != nil {
		return
	}
	data = &proto.Data{
		Bytes: fileBytes,
	}

	return data, err
}

// Processes implements the proto.OSDServer interface.
func (r *Registrator) Processes(ctx context.Context, in *empty.Empty) (reply *proto.ProcessesReply, err error) {
	ctx = namespaces.WithNamespace(ctx, "system")
	client, err := containerd.New(constants.ContainerdSocket)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	containers, err := client.Containers(ctx)
	processes := []*proto.Process{}
	for _, c := range containers {
		info, err := c.Info(ctx)
		if err != nil {
			log.Println(err)
			continue
		}

		task, err := c.Task(ctx, nil)
		if err != nil {
			log.Println(err)
			continue
		}

		status, err := task.Status(ctx)
		if err != nil {
			log.Println(err)
			continue
		}

		process := &proto.Process{
			Id:     task.ID(),
			Image:  info.Image,
			Status: string(status.Status),
		}

		if status.Status == containerd.Running {
			metrics, err := task.Metrics(ctx)
			if err != nil {
				log.Println(err)
				continue
			}
			anydata, err := typeurl.UnmarshalAny(metrics.Data)
			if err != nil {
				log.Println(err)
				continue
			}
			data, ok := anydata.(*cgroups.Metrics)
			if !ok {
				log.Println(errors.New("failed to convert metric data to cgroups.Metrics"))
				continue
			}
			process.MemoryUsage = data.Memory.Usage.Usage
			process.CpuUsage = data.CPU.Usage.Total
		}

		processes = append(processes, process)
	}

	reply = &proto.ProcessesReply{Processes: processes}

	return reply, nil
}

// Restart implements the proto.OSDServer interface.
func (r *Registrator) Restart(ctx context.Context, in *proto.RestartRequest) (reply *proto.RestartReply, err error) {
	ctx = namespaces.WithNamespace(ctx, "system")
	client, err := containerd.New(constants.ContainerdSocket)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	task := client.TaskService()
	_, err = task.Kill(ctx, &tasks.KillRequest{ContainerID: in.Id, Signal: uint32(unix.SIGTERM)})
	if err != nil {
		return nil, err
	}

	reply = &proto.RestartReply{}

	return
}

// Reset implements the proto.OSDServer interface.
func (r *Registrator) Reset(ctx context.Context, in *empty.Empty) (reply *proto.ResetReply, err error) {
	// TODO(andrewrynhard): Delete all system tasks and containers.

	// Set the process arguments.
	args := runner.Args{
		ID:          "reset",
		ProcessArgs: []string{"/bin/kubeadm", "reset", "--force"},
	}

	// Set the mounts.
	// nolint: dupl
	mounts := []specs.Mount{
		{Type: "cgroup", Destination: "/sys/fs/cgroup", Options: []string{"ro"}},
		{Type: "bind", Destination: "/var/run", Source: "/run", Options: []string{"rbind", "rshared", "rw"}},
		{Type: "bind", Destination: "/var/lib/docker", Source: "/var/lib/docker", Options: []string{"rbind", "rshared", "rw"}},
		{Type: "bind", Destination: "/var/lib/kubelet", Source: "/var/lib/kubelet", Options: []string{"rbind", "rshared", "rw"}},
		{Type: "bind", Destination: "/etc/kubernetes", Source: "/var/etc/kubernetes", Options: []string{"bind", "rw"}},
		{Type: "bind", Destination: "/etc/os-release", Source: "/etc/os-release", Options: []string{"bind", "ro"}},
		{Type: "bind", Destination: "/bin/crictl", Source: "/bin/crictl", Options: []string{"bind", "ro"}},
		{Type: "bind", Destination: "/bin/kubeadm", Source: "/bin/kubeadm", Options: []string{"bind", "ro"}},
		{Type: "bind", Destination: "/bin/kubeadm.sh", Source: "/run/kubeadm.sh", Options: []string{"bind", "ro"}},
	}

	cr := containerdrunner.Containerd{}

	err = cr.Run(
		r.Data,
		args,
		runner.WithContainerImage(constants.KubernetesImage),
		runner.WithOCISpecOpts(
			containerdrunner.WithMemoryLimit(int64(1000000*512)),
			containerdrunner.WithRootfsPropagation("slave"),
			oci.WithMounts(mounts),
			oci.WithHostNamespace(specs.PIDNamespace),
			oci.WithParentCgroupDevices,
			oci.WithPrivileged,
		),
		runner.WithType(runner.Once),
	)

	if err != nil {
		return nil, err
	}

	reply = &proto.ResetReply{}

	return reply, nil
}

// Reboot implements the proto.OSDServer interface.
func (r *Registrator) Reboot(ctx context.Context, in *empty.Empty) (reply *proto.RebootReply, err error) {
	unix.Reboot(int(unix.LINUX_REBOOT_CMD_RESTART))

	reply = &proto.RebootReply{}

	return
}

// Dmesg implements the proto.OSDServer interface. The klogctl syscall is used
// to read from the ring buffer at /proc/kmsg by taking the
// SYSLOG_ACTION_READ_ALL action. This action reads all messages remaining in
// the ring buffer non-destructively.
func (r *Registrator) Dmesg(ctx context.Context, in *empty.Empty) (data *proto.Data, err error) {
	// Return the size of the kernel ring buffer
	size, err := unix.Klogctl(constants.SYSLOG_ACTION_SIZE_BUFFER, nil)
	if err != nil {
		return
	}
	// Read all messages from the log (non-destructively)
	buf := make([]byte, size)
	n, err := unix.Klogctl(constants.SYSLOG_ACTION_READ_ALL, buf)
	if err != nil {
		return
	}

	data = &proto.Data{Bytes: buf[:n]}

	return data, err
}

// Logs implements the proto.OSDServer interface. Service or container logs can
// be requested and the contents of the log file are streamed in chunks.
func (r *Registrator) Logs(req *proto.LogsRequest, l proto.OSD_LogsServer) (err error) {
	var chunk chunker.ChunkReader
	if req.Container {
		switch r.Data.Services.Init.ContainerRuntime {
		case constants.ContainerRuntimeDocker:
			chunk, err = dockerLogs(req.Process)
			if err != nil {
				return
			}
		case constants.ContainerRuntimeCRIO:
			chunk, err = crioLogs(req.Process)
			if err != nil {
				return
			}
		}
	} else {
		logpath := servicelog.FormatLogPath(req.Process)
		file, _err := os.OpenFile(logpath, os.O_RDONLY, 0)
		if _err != nil {
			err = _err
			return
		}
		chunk = filechunker.NewChunker(file)
	}

	if chunk == nil {
		return fmt.Errorf("no log reader found")
	}

	for data := range chunk.Read(l.Context()) {
		if err = l.Send(&proto.Data{Bytes: data}); err != nil {
			return
		}
	}

	return nil
}

// Version implements the proto.OSDServer interface.
func (r *Registrator) Version(ctx context.Context, in *empty.Empty) (data *proto.Data, err error) {
	v, err := version.NewVersion()
	if err != nil {
		return
	}

	data = &proto.Data{Bytes: []byte(v)}

	return data, err
}

func crioLogs(id string) (chunk chunker.Chunker, err error) {
	cli, err := crioclient.New(constants.ContainerRuntimeCRIOSocket)
	if err != nil {
		return
	}
	info, err := cli.ContainerInfo(id)
	if err != nil {
		return
	}
	file, err := os.OpenFile(info.LogPath, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	chunk = filechunker.NewChunker(file)

	return chunk, nil
}

func dockerLogs(id string) (chunk chunker.Chunker, err error) {
	cli, err := dockerclient.NewEnvClient()
	if err != nil {
		return
	}
	stream, err := cli.ContainerLogs(context.Background(), id, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: false,
		Follow:     true,
		Tail:       "40",
	})
	if err != nil {
		return
	}

	chunk = streamchunker.NewChunker(stream)

	return chunk, nil
}

package connection

import (
	"bytes"
	"fmt"
	"net/netip"

	"golang.org/x/crypto/ssh"
)

// SSHConnection encapsulates the SSH connection to the devices as well as any commands we run on them
type SSHConnection struct {
	connection *ssh.Client
}

// NewSSHConnection returns a new SSH connection struct
func NewSSHConnection(ip netip.Addr, port uint16, user, password string) (*SSHConnection, error) {
	// SSH client configuration
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		// Note: The following is insecure for production use. It skips SSH key verification.
		// For production code, replace it with a callback that verifies the server's identity.
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connecting to the SSH server
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip.String(), port), config)
	if err != nil {
		return nil, err
	}

	return &SSHConnection{connection: connection}, nil
}

func (s *SSHConnection) remoteCommand(command string) (*bytes.Buffer, error) {
	session, err := s.connection.NewSession()
	if err != nil {
		return nil, err
	}
	defer func(session *ssh.Session) {
		_ = session.Close()
	}(session)

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(command); err != nil {
		return &b, err
	}

	return &b, nil
}

// FetchLiveConfig gets the current config from the host
func (s *SSHConnection) FetchLiveConfig() ([]byte, error) {
	buf, err := s.remoteCommand("cat /config/config.boot")
	if err != nil {
		return nil, fmt.Errorf("error reading existing config: %w", err)
	}

	return buf.Bytes(), nil
}

// WriteFile writes a file to the remote host
func (s *SSHConnection) WriteFile(remotePath string, contents []byte) error {
	buf, err := s.remoteCommand(fmt.Sprintf("echo '%s' > %s", string(contents), remotePath))
	fmt.Print(buf.String())
	return err
}

// DeleteFile deletes a file on the remote host
func (s *SSHConnection) DeleteFile(remotePath string) error {
	buf, err := s.remoteCommand(fmt.Sprintf("rm %s", remotePath))
	fmt.Print(buf.String())
	return err
}

// ApplyConfig applies, commits, and saves the config at the supplied path
func (s *SSHConnection) ApplyConfig(configPath string) error {
	commands := []string{
		"/opt/vyatta/sbin/vyatta-cfg-cmd-wrapper begin",
		fmt.Sprintf("/opt/vyatta/sbin/vyatta-cfg-cmd-wrapper load %s", configPath),
		"/opt/vyatta/sbin/vyatta-cfg-cmd-wrapper commit",
		"/opt/vyatta/sbin/vyatta-cfg-cmd-wrapper save",
	}

	for _, cmd := range commands {
		buf, err := s.remoteCommand(cmd)
		if buf != nil {
			fmt.Print(buf.String())
		}
		if err != nil {
			return fmt.Errorf("error running command %s: %w", cmd, err)
		}
	}

	return nil
}

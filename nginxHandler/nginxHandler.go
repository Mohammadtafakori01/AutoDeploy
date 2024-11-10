package nginxHandler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const nginxConfigTemplate = `
server {
    listen 80;
    server_name {{.URL}};
    root {{.RootDir}};
    index index.html index.htm;
    location / {
        try_files $uri $uri/ =404;
    }
}
`

// CreateNginxConfig generates an Nginx configuration file for the given URL and root directory.
func CreateNginxConfig(url, rootDir, configPath string) error {
	// Prepare data for the template
	data := struct {
		URL     string
		RootDir string
	}{
		URL:     url,
		RootDir: rootDir,
	}

	// Parse and execute the template
	tmpl, err := template.New("nginxConfig").Parse(nginxConfigTemplate)
	if err != nil {
		return err
	}

	// Create the configuration file
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return err
	}

	// Create a symlink in sites-enabled
	symlinkPath := "/etc/nginx/sites-enabled/" + filepath.Base(configPath)
	if _, err := os.Lstat(symlinkPath); os.IsNotExist(err) {
		if err := os.Symlink(configPath, symlinkPath); err != nil {
			return err
		}
	}

	return nil
}

// ReloadNginx reloads the Nginx service to apply new configurations.
func ReloadNginx() error {
	cmd := exec.Command("nginx", "-s", "reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to reload nginx: %v, output: %s", err, string(output))
	}
	return nil
}

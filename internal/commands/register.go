package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type RegisterCmd struct{}

func (c *RegisterCmd) Run(args []string) error {
	fmt.Fprintln(os.Stderr, "Basecamp OAuth App Registration Helper")
	fmt.Fprintln(os.Stderr, strings.Repeat("=", 40))
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "This helper will generate the values you need to register")
	fmt.Fprintln(os.Stderr, "your Basecamp OAuth application.")
	fmt.Fprintln(os.Stderr)

	reader := bufio.NewReader(os.Stdin)

	appName := prompt(reader, "Application name", "My Basecamp CLI")
	companyName := prompt(reader, "Company/Organization name", "")
	websiteURL := prompt(reader, "Website URL", "https://github.com/robzolkos/basecamp-cli")
	accessibleURL := prompt(reader, "URL where this computer is accessible (e.g., https://myhost.tailscale.ts.net)", "")

	// Build redirect URI from accessible URL
	redirectURI := buildRedirectURI(accessibleURL)

	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, strings.Repeat("=", 60))
	fmt.Fprintln(os.Stderr, "REGISTRATION INSTRUCTIONS")
	fmt.Fprintln(os.Stderr, strings.Repeat("=", 60))
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "1. Visit: https://launchpad.37signals.com/integrations")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "2. Click 'Register another application'")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "3. Fill out the form with these values:")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "   Name of your application:  %s\n", appName)
	fmt.Fprintf(os.Stderr, "   Your company/organization: %s\n", companyName)
	fmt.Fprintf(os.Stderr, "   Website URL:               %s\n", websiteURL)
	fmt.Fprintf(os.Stderr, "   Redirect URI:              %s\n", redirectURI)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "4. After registering, copy your Client ID and Client Secret")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "5. Run 'basecamp init' and enter the credentials when prompted")
	fmt.Fprintln(os.Stderr, "   (use the same Redirect URI shown above)")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "6. Run 'basecamp auth' to authenticate")
	fmt.Fprintln(os.Stderr, strings.Repeat("=", 60))

	return PrintJSON(map[string]string{
		"application_name": appName,
		"company_name":     companyName,
		"website_url":      websiteURL,
		"redirect_uri":     redirectURI,
		"registration_url": "https://launchpad.37signals.com/integrations",
	})
}

func buildRedirectURI(accessibleURL string) string {
	if accessibleURL == "" {
		return "http://localhost:3002/callback"
	}

	// Remove trailing slash if present
	accessibleURL = strings.TrimSuffix(accessibleURL, "/")

	// Add port and callback path
	return accessibleURL + ":3002/callback"
}

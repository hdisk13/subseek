package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/subscriptions"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func readCredentialsFromFile() (string, string, string, error) {
	file, err := os.Open("creds.config")
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	clientID, clientSecret, tenantID := "", "", ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "AZURE_CLIENT_ID=") {
			clientID = strings.TrimPrefix(line, "AZURE_CLIENT_ID=")
		} else if strings.HasPrefix(line, "AZURE_CLIENT_SECRET=") {
			clientSecret = strings.TrimPrefix(line, "AZURE_CLIENT_SECRET=")
		} else if strings.HasPrefix(line, "AZURE_TENANT_ID=") {
			tenantID = strings.TrimPrefix(line, "AZURE_TENANT_ID=")
		}
	}

	return clientID, clientSecret, tenantID, scanner.Err()
}

func main() {
	ctx := context.Background()

	// Configure the Azure SDK for Go with your Azure credentials
	authorizer, err := auth.NewClientCredentialsConfig("<YOUR_AZURE_CLIENT_ID>", "<YOUR_AZURE_CLIENT_SECRET>", "<YOUR_AZURE_TENANT_ID>").Authorizer()
	if err != nil {
		fmt.Printf("Failed to get Azure authorizer: %v\n", err)
		return
	}

	client := subscriptions.NewClient()
	client.Authorizer = authorizer

	// List Azure subscriptions
	subs, err := client.List(ctx)
	if err != nil {
		fmt.Printf("Failed to list Azure subscriptions: %v\n", err)
		return
	}

	fmt.Printf("%-40s %-36s %-36s %-10s\n", "Subscription_Name", "Subscription", "Tenant", "Is_Default")
	for _, sub := range subs.Values() {
		fmt.Printf("%-40s %-36s %-36s %-10v\n", *sub.DisplayName, *sub.SubscriptionID, *sub.TenantID, sub.State == "Enabled")
	}

	fmt.Println("\n-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-\n")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Activate which subscription?  (enter==don't change) >> ")
	newSub, _ := reader.ReadString('\n')

	if newSub == "" || newSub == "\n" {
		fmt.Println("Keeping current subscription.")
	} else {
		_, err := client.Set(ctx, newSub[:len(newSub)-1])
		if err != nil {
			fmt.Printf("Failed to set Azure subscription: %v\n", err)
			return
		}
		fmt.Println("Subscription changed.")
	}
}

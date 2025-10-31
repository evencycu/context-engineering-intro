package com.teamsnotify.sdk.example;

import org.openapitools.client.*;
import org.openapitools.client.api.*;
import org.openapitools.client.model.*;
import java.util.*;

public class Example {
    public static void main(String[] args) {
        // Initialize the client
        ApiClient client = Configuration.getDefaultApiClient();
        client.setBasePath("http://localhost:8080");

        ExternalApiApi externalApi = new ExternalApiApi(client);

        try {
            // Send a notification using external API
            ExternalNotifyRequest request = new ExternalNotifyRequest();
            request.setNotifyKey("cfh-alert-gogo");
            request.setMessage("Hello from Java SDK!");
            request.setMessageType("text");
            request.setPriority("normal");
            request.setTargets(Arrays.asList("all"));

            ExternalNotifyResponse response = externalApi.apiV1NotifyPost(request);
            System.out.println("Notification sent: " + response.getNotificationId());

            // Get project destinations
            ProjectDestinationsResponse dests = externalApi.apiV1DestinationsNotifyKeyGet("cfh-alert-gogo");
            System.out.println("Found " + dests.getDestinations().size() + " destinations");

        } catch (ApiException e) {
            System.err.println("Error: " + e.getMessage());
            e.printStackTrace();
        }
    }
}

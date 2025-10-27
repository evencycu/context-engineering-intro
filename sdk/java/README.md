# Teams Notification Java SDK

Java SDK for Teams Notification API

## Installation

### Maven

```xml
<dependencies>
    <dependency>
        <groupId>com.teamsnotify</groupId>
        <artifactId>teams-notify-sdk</artifactId>
        <version>1.0.0</version>
    </dependency>
</dependencies>
```

### Gradle

```gradle
dependencies {
    implementation 'com.teamsnotify:teams-notify-sdk:1.0.0'
}
```

## Usage

### Basic Example

```java
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
            request.setNotifyKey("your-notify-key");
            request.setMessage("Hello from Java SDK!");
            request.setMessageType("text");
            request.setPriority("normal");
            request.setTargets(Arrays.asList("all"));
            
            ExternalNotifyResponse response = externalApi.apiV1NotifyPost(request);
            System.out.println("Notification sent: " + response.getNotificationId());
            
            // Get project destinations
            ProjectDestinationsResponse dests = externalApi.apiV1DestinationsNotifyKeyGet("your-notify-key");
            System.out.println("Found " + dests.getDestinations().size() + " destinations");
            
        } catch (ApiException e) {
            System.err.println("Error: " + e.getMessage());
            e.printStackTrace();
        }
    }
}
```

### External API Example

```java
// Send notification via external API
ExternalNotifyRequest request = new ExternalNotifyRequest();
request.setNotifyKey("your-notify-key");
request.setMessage("Hello from external API!");
request.setTargets(Arrays.asList("all"));

ExternalNotifyResponse response = externalApi.apiV1NotifyPost(request);

// Get project destinations
ProjectDestinationsResponse dests = externalApi.apiV1DestinationsNotifyKeyGet("your-notify-key");
```

### Internal API Example

```java
// Create authenticated client
ApiClient client = Configuration.getDefaultApiClient();
client.setBasePath("http://localhost:8080");
client.setApiKey("your-api-key");

// Use protected endpoints
InternalApiUsersApi usersApi = new InternalApiUsersApi(client);
InternalApiProjectsApi projectsApi = new InternalApiProjectsApi(client);

List<User> users = usersApi.internalV1UsersGet();
List<Project> projects = projectsApi.internalV1ProjectsGet();
```

## API Reference

### External API
- `apiV1NotifyPost` - Send notification
- `apiV1DestinationsNotifyKeyGet` - Get project destinations

### Internal API
- Users: `internalV1UsersGet`, `internalV1UsersPost`, etc.
- Projects: `internalV1ProjectsGet`, `internalV1ProjectsPost`, etc.
- Notifications: `internalV1NotificationsGet`, `internalV1NotificationsPost`, etc.
- Destinations: `internalV1DestinationsGet`, `internalV1DestinationsPost`, etc.
- Bots: `internalV1BotsPlatformGet`, `internalV1BotsPlatformPost`, etc.

## License

MIT License
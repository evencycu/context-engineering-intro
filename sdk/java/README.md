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
import com.teamsnotify.sdk.Client;
import com.teamsnotify.sdk.models.*;
import com.teamsnotify.sdk.services.*;

public class Example {
    public static void main(String[] args) {
        // Initialize the client
        Client client = new Client.Builder()
            .baseUrl("http://localhost:8080")
            .build();
        
        // Send a notification
        SendNotificationRequest request = SendNotificationRequest.builder()
            .notifyKey("your-notify-key")
            .message("Hello from Java SDK!")
            .messageType("text")
            .priority("normal")
            .targets(Arrays.asList("all"))
            .build();
        
        try {
            SendNotificationResponse response = client.notifications().send(request);
            System.out.println("Notification sent: " + response.getNotificationId());
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
```

### External API Example

```java
// Use external API (no authentication required)
ExternalClient client = new ExternalClient.Builder()
    .baseUrl("http://localhost:8080/api/v1")
    .build();

// Send notification via external API
SendNotificationRequest request = SendNotificationRequest.builder()
    .notifyKey("your-notify-key")
    .message("Hello from external API!")
    .targets(Arrays.asList("all"))
    .build();

SendNotificationResponse response = client.sendNotification(request);

// Get project destinations
ProjectDestinationsResponse dests = client.getProjectDestinations("your-notify-key");
```

### Authentication

```java
// Create authenticated client
Client client = new Client.Builder()
    .baseUrl("http://localhost:8080")
    .apiKey("your-api-key")
    .build();

// Use protected endpoints
List<User> users = client.users().list();
List<Project> projects = client.projects().list();
```

## API Reference

### External API

- `SendNotification(SendNotificationRequest request)`
- `GetProjectDestinations(String notifyKey)`

### Internal API

#### Users
- `List<List<User>>()`
- `Get(String id)`
- `Create(CreateUserRequest request)`
- `Update(String id, UpdateUserRequest request)`
- `Delete(String id)`

#### Projects
- `List<List<Project>>()`
- `Get(String id)`
- `Create(CreateProjectRequest request)`
- `Update(String id, UpdateProjectRequest request)`
- `Delete(String id)`

#### Notifications
- `Send(SendNotificationRequest request)`
- `Get(String id)`
- `List(NotificationQuery query)`
- `Cancel(String id)`

## License

MIT License


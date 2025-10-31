import java.io.*;
import java.net.*;
import java.nio.charset.StandardCharsets;
import java.util.*;

public class SimpleJavaTest {
    public static void main(String[] args) {
        String baseURL = "http://localhost:8080";

        // Test 1: Send notification
        System.out.println("Testing external API - Send notification...");

        try {
            String requestBody = "{\n" +
                    "  \"notifyKey\": \"cfh-alert-gogo\",\n" +
                    "  \"message\": \"Hello from Java SDK test!\",\n" +
                    "  \"messageType\": \"text\",\n" +
                    "  \"priority\": \"normal\",\n" +
                    "  \"targets\": [\"all\"]\n" +
                    "}";

            URL url = new URL(baseURL + "/api/v1/notify");
            HttpURLConnection conn = (HttpURLConnection) url.openConnection();
            conn.setRequestMethod("POST");
            conn.setRequestProperty("Content-Type", "application/json");
            conn.setDoOutput(true);

            try (OutputStream os = conn.getOutputStream()) {
                byte[] input = requestBody.getBytes(StandardCharsets.UTF_8);
                os.write(input, 0, input.length);
            }

            int responseCode = conn.getResponseCode();
            String response = readResponse(conn);

            if (responseCode == 200) {
                System.out.println("✅ Notification sent successfully!");
                System.out.println("   Response: " + response);
            } else {
                System.out.println("❌ Notification failed with status " + responseCode);
                System.out.println("   Response: " + response);
            }

        } catch (Exception e) {
            System.err.println("Error sending notification: " + e.getMessage());
        }

        // Test 2: Get project destinations
        System.out.println("\nTesting external API - Get project destinations...");

        try {
            URL url = new URL(baseURL + "/api/v1/destinations/cfh-alert-gogo");
            HttpURLConnection conn = (HttpURLConnection) url.openConnection();
            conn.setRequestMethod("GET");

            int responseCode = conn.getResponseCode();
            String response = readResponse(conn);

            if (responseCode == 200) {
                System.out.println("✅ Project destinations retrieved successfully!");
                System.out.println("   Response: " + response);
            } else {
                System.out.println("❌ Get destinations failed with status " + responseCode);
                System.out.println("   Response: " + response);
            }

        } catch (Exception e) {
            System.err.println("Error getting destinations: " + e.getMessage());
        }

        System.out.println("\n🎉 All external API tests completed!");
    }

    private static String readResponse(HttpURLConnection conn) throws IOException {
        InputStream inputStream;
        if (conn.getResponseCode() >= 200 && conn.getResponseCode() < 300) {
            inputStream = conn.getInputStream();
        } else {
            inputStream = conn.getErrorStream();
        }

        StringBuilder response = new StringBuilder();
        try (BufferedReader br = new BufferedReader(new InputStreamReader(inputStream))) {
            String line;
            while ((line = br.readLine()) != null) {
                response.append(line);
            }
        }
        return response.toString();
    }
}

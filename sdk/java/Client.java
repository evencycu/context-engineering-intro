package com.teamsnotify.sdk;

import com.teamsnotify.sdk.services.*;

public class Client {
    private final String baseUrl;
    private final String apiKey;
    private final OkHttpClient httpClient;
    
    private Client(Builder builder) {
        this.baseUrl = builder.baseUrl;
        this.apiKey = builder.apiKey;
        this.httpClient = new OkHttpClient();
    }
    
    public static class Builder {
        private String baseUrl;
        private String apiKey;
        
        public Builder baseUrl(String baseUrl) {
            this.baseUrl = baseUrl;
            return this;
        }
        
        public Builder apiKey(String apiKey) {
            this.apiKey = apiKey;
            return this;
        }
        
        public Client build() {
            return new Client(this);
        }
    }
    
    public NotificationService notifications() {
        return new NotificationService(this);
    }
    
    public UserService users() {
        return new UserService(this);
    }
    
    public ProjectService projects() {
        return new ProjectService(this);
    }
    
    // ... other services
}

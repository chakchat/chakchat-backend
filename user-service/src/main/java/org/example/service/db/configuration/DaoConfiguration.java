package org.example.service.db.configuration;

import org.example.service.db.UserDAO;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class DaoConfiguration {

    @Bean
    UserDAO userDAO() {
        return new UserDAO();
    }
}

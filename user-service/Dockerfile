# Используем образ с Java 23
FROM eclipse-temurin:23-jdk AS build

# Установка рабочего каталога
WORKDIR /app

# Копируем файлы Gradle (чтобы кэшировать зависимости)
COPY gradle gradle
COPY build.gradle settings.gradle ./

# Копируем исходный код
COPY src ./src

# Собираем проект
RUN ./gradlew build --no-daemon

# Создаем образ для запуска
FROM eclipse-temurin:23-jre

# Установка рабочего каталога
WORKDIR /app

# Копируем скомпилированный JAR-файл из предыдущего образа
COPY --from=build /app/build/libs/*.jar app.jar

# Экспонируем порт
EXPOSE 8080

# Определяем команду для запуска приложения
ENTRYPOINT ["java", "-jar", "app.jar"]
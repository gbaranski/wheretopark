FROM gradle:jdk17 AS build
COPY --chown=gradle:gradle . /home/gradle/src
WORKDIR /home/gradle/src
RUN gradle storekeeper:buildFatJar --no-daemon

FROM openjdk:17
RUN mkdir /app
COPY --from=build /home/gradle/src/storekeeper/build/libs/storekeeper-all.jar /app/storekeeper.jar
ENTRYPOINT ["java","-jar","/app/storekeeper.jar"]

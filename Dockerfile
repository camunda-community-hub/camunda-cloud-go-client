FROM scratch
EXPOSE 8080
ENTRYPOINT ["/camunda-cloud-go-client"]
COPY ./bin/ /
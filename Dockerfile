FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-vgs"]
COPY baton-vgs /
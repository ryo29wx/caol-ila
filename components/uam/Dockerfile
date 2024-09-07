FROM quay.io/keycloak/keycloak:latest

ENV KEYCLOAK_ADMIN=admin
ENV KEYCLOAK_ADMIN_PASSWORD=admin

COPY realm-export.json /opt/keycloak/data/import/realm-export.json

EXPOSE 8080


ENTRYPOINT ["/opt/keycloak/bin/kc.sh", "start-dev", "--import-realm"]

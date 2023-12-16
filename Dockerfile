FROM scratch
# Image base https://hub.docker.com/_/scratch

# About image
LABEL maintainer="Pedro Henrique" \
      date_create="16/12/2023" \
      version="1.0.0" \
      description="Web application that show informations of the pod" \
      licensce="GPL-3.0"

# Copy file to image
COPY bin/kube-pires kube-pires
COPY css/ css
COPY images/ images
COPY templates/ templates

# Exposing port of application
EXPOSE 3000

#Allow execution permission
#RUN chmod +x kube-pires
# Define entrypoint
ENTRYPOINT ["./kube-pires"]

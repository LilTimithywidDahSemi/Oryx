# Use the mariner base image
FROM mcr.microsoft.com/dotnet/sdk:7.0-cbl-mariner2.0 as main

ARG SDK_STORAGE_BASE_URL_VALUE="https://oryx-cdn.microsoft.io"
ARG AI_CONNECTION_STRING

COPY --from=oryxdevmcr.azurecr.io/private/oryx/buildscriptgenerator /opt/buildscriptgen/ /opt/buildscriptgen/
COPY --from=oryxdevmcr.azurecr.io/private/oryx/support-files-image-for-build /tmp/oryx/ /opt/tmp

# Install an assortment of traditional tooling (unicode, SSL, HTTP, etc.)
# Also install some PHP requirements  
RUN tdnf update -y \
    && tdnf install -y \
# .NET Core dependencies for running Oryx
        icu \ 
        gcc \ 
      #  rsync \
      #  libgdiplus \
         # Required for mysqlclient
        mariadb \
       # mysql \
    && chmod a+x /opt/buildscriptgen/GenerateBuildScript \
    && mkdir -p /opt/oryx \
    && ln -s /opt/buildscriptgen/GenerateBuildScript /opt/oryx/oryx \
    # Install Python SDKs
    # Upgrade system python
    && PYTHONIOENCODING="UTF-8" \
    && tmpDir="/opt/tmp" \
    && mkdir -p /usr/local/share/pip-cache/lib \
    && chmod -R 777 /usr/local/share/pip-cache \
    && mkdir -p /opt/dotnet/7.0.202/shared/Microsoft.NETCore.App/7.0.4 \
    && chmod -R 777 /opt/dotnet/7.0.202/shared/Microsoft.NETCore.App/7.0.4 \
   # && cp -R /usr/lib/python3.9 /opt/python/3.9.16/bin \
   # && ln -s opt/python/3.9.16/bin/python3.9 /usr/lib/python3.9 \
    && echo "dotnet-mariner" > /opt/oryx/.imagetype \
    && echo "MARINER" | tr '[a-z]' '[A-Z]' > /opt/oryx/.ostype

RUN tmpDir="/opt/tmp" \
    && cp -f $tmpDir/images/build/benv.sh /opt/oryx/benv \
    && cp -f $tmpDir/images/build/logger.sh /opt/oryx/logger \
    && chmod +x /opt/oryx/benv \
    && chmod +x /opt/oryx/logger

ENV ORYX_SDK_STORAGE_BASE_URL=${SDK_STORAGE_BASE_URL_VALUE} \
    ENABLE_DYNAMIC_INSTALL="false" \
    DYNAMIC_INSTALL_ROOT_DIR="/opt" \
    PYTHONIOENCODING="UTF-8" \
    LANG="C.UTF-8" \
    LANGUAGE="C.UTF-8" \
    LC_ALL="C.UTF-8" \
    PATH="/usr/local/go/bin:/opt/python/latest/bin:/opt/oryx:/opt/yarn/stable/bin:/opt/hugo/lts:/opt/dotnet:$PATH" \
    ORYX_AI_CONNECTION_STRING="${AI_CONNECTION_STRING}" \
    DOTNET_SKIP_FIRST_TIME_EXPERIENCE="1"

ENTRYPOINT [ "benv" ] 
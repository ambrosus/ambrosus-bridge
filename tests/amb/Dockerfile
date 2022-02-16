FROM openethereum/openethereum:latest

USER root

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    apk add jq

ARG NETWORK_ID
ARG NETWORK_NAME
ARG USER_ID
ARG GROUP_ID
ARG USER_NAME
ARG GROUP_NAME
ARG BASE_PATH
ARG ACCOUNT_PASSWORD
ARG REWRITE_CONFIG_FILES

ARG user=${USER_NAME}
ARG group=${GROUP_NAME}
ARG uid=${USER_ID}
ARG gid=${GROUP_ID}

ENV BASE_PATH ${BASE_PATH}

WORKDIR ${BASE_PATH}

RUN mkdir -p ${BASE_PATH}/data
RUN mkdir -p ${BASE_PATH}/config

COPY ./config/ ${BASE_PATH}/config/
COPY ./configuration.sh ${BASE_PATH}/configuration.sh

RUN sh ${BASE_PATH}/configuration.sh ${REWRITE_CONFIG_FILES} ${BASE_PATH} ${NETWORK_ID} ${NETWORK_NAME}

RUN chown -R ${USER_NAME}:${GROUP_NAME} ${BASE_PATH}
RUN chmod a+rwx -R ${BASE_PATH}

USER ${USER_NAME}

EXPOSE 8545 8546 30303

ENTRYPOINT ./openethereum \
    --config ${BASE_PATH}/config/config.toml \
    --base-path=${BASE_PATH}/data \

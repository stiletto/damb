ARG damb_prefix
ARG damb_tag
FROM ${damb_prefix}hello-c.build:${damb_tag} AS build
FROM ${damb_prefix}hello-c.base:${damb_tag}
COPY --from=build /root/hello /bin/hello-c
ENTRYPOINT ["/bin/hello-c"]

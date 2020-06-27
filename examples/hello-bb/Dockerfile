# DAMB:ctx: tar -cv hello-bb
ARG damb_prefix
ARG damb_tag
FROM ${damb_prefix}hello-go:${damb_tag}

ADD hello-bb/hello.sh /bin/hello.sh

ENTRYPOINT ["/bin/hello.sh"]

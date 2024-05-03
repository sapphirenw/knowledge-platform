FROM pgvector/pgvector:pg16

# update package lists
RUN apt-get update
RUN apt install -y curl

# copy over the postgres config file
COPY ./postgresql.conf /etc/postgresql.conf

# Set the environment variable to store the directory
# where initialization scripts are to be executed
ENV INIT_DIR /docker-entrypoint-initdb.d

# Create the directory if it doesn't exist
RUN mkdir -p $INIT_DIR

# Set the working directory to the init directory
WORKDIR $INIT_DIR

# Copy initialization scripts to the container
COPY ./schema/00_GENERATED_schema.sql init.sql

# Grant permissions on the scripts to be executable
RUN chmod +x $INIT_DIR/*

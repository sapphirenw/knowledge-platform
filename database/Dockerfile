FROM pgvector/pgvector:pg16

# update package lists
RUN apt-get update
RUN apt install -y curl

# copy over the postgres config file
COPY ./postgresql.conf /etc/postgresql.conf

# Set the environment variable to store the directory
# where initialization scripts are to be executed
ENV INIT_DIR /docker-entrypoint-initdb.d

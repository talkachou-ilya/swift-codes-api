FROM mongo:8

WORKDIR /data

COPY ./swiftcodes.json .

CMD mongoimport --host mongo \
    --port 27017 \
    --db "${MONGO_NAME}" \
    --username "${MONGO_USERNAME}" \
    --password "${MONGO_PASSWORD}" \
    --authenticationDatabase admin \
    --collection swift-codes \
    --file swiftcodes.json \
    --jsonArray \
    --drop

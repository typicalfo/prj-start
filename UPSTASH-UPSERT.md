# Upsert Data

> Upserts (inserts or updates) the raw text data after embedding it.

<Warning>
  To use this endpoint, the index must be created with an [embedding model](/vector/features/embeddingmodels).
</Warning>

<Tip>
  Vector embedding of the raw text data will be upserted into the
  default namespace by default.
  You can use a different namespace by specifying it in the request path.
</Tip>

## Request

You can either upsert a single data, or multiple data in an array.

<ParamField body="id" type="string" required>
  The id of the vector.
</ParamField>

<ParamField body="data" type="string" required>
  The raw text data to embed and upsert.
</ParamField>

<ParamField body="metadata" type="Object">
  The metadata of the vector. This makes identifying vectors
  on retrieval easier and can be used to with filters on queries.
</ParamField>

<Note>
  Data field of the vector will be automatically set to the
  raw text data, so that you can access it later, during
  queries.
</Note>

## Path

<ParamField path="namespace" type="string" default="">
  The namespace to use.
  When no namespace is specified, the default namespace will be used.
</ParamField>

## Response

<ResponseField name="result" type="string">
  `"Success"` string.
</ResponseField>

<RequestExample>
  ```sh curl theme={"system"}
  curl $UPSTASH_VECTOR_REST_URL/upsert-data \
    -X POST \
    -H "Authorization: Bearer $UPSTASH_VECTOR_REST_TOKEN" \
    -d '[ 
      { "id": "id-0", "data": "Upstash is a serverless data platform.", "metadata": { "link": "upstash.com" } }, 
      { "id": "id-1", "data": "Upstash Vector is a serverless vector database." }
    ]'
  ```

  ```sh curl (Namespace) theme={"system"}
  curl $UPSTASH_VECTOR_REST_URL/upsert-data/ns \
    -X POST \
    -H "Authorization: Bearer $UPSTASH_VECTOR_REST_TOKEN" \
    -d '{ "id": "id-2", "data": "Upstash is a serverless data platform.", "metadata": { "link": "upstash.com" } }'
  ```
</RequestExample>

<ResponseExample>
  ```json 200 OK theme={"system"}
  {
      "result": "Success"
  }
  ```

  ```json 422 Unprocessable Entity theme={"system"}
  {
      "error": "Embedding data for this index is not allowed. The index must be created with an embedding model to use it.",
      "status": 422
  }
  ```
</ResponseExample>

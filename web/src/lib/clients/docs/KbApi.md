# KbApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**deleteDocument**](KbApi.md#deletedocument) | **DELETE** /kb/{memory}/documents/{document} |  |
| [**listDocuments**](KbApi.md#listdocuments) | **GET** /kb/{memory}/documents |  |
| [**listKbMemories**](KbApi.md#listkbmemories) | **GET** /kb/_memories |  |
| [**searchKb**](KbApi.md#searchkb) | **GET** /kb/{memory} |  |
| [**submitDocument**](KbApi.md#submitdocument) | **POST** /kb/{memory}/documents/{document} |  |



## deleteDocument

> deleteDocument(memory, document)



deletes all the knowledge chunks identified by a specific file name

### Example

```ts
import {
  Configuration,
  KbApi,
} from '';
import type { DeleteDocumentRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new KbApi(config);

  const body = {
    // string
    memory: memory_example,
    // string
    document: document_example,
  } satisfies DeleteDocumentRequest;

  try {
    const data = await api.deleteDocument(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **memory** | `string` |  | [Defaults to `undefined`] |
| **document** | `string` |  | [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | knowledge is deleted |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listDocuments

> Array&lt;string&gt; listDocuments(memory)



lists all the documents in a memory slot

### Example

```ts
import {
  Configuration,
  KbApi,
} from '';
import type { ListDocumentsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new KbApi(config);

  const body = {
    // string
    memory: memory_example,
  } satisfies ListDocumentsRequest;

  try {
    const data = await api.listDocuments(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **memory** | `string` |  | [Defaults to `undefined`] |

### Return type

**Array<string>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | documents are returned |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listKbMemories

> { [key: string]: Memory; } listKbMemories()



lists all the knowledge base memory slots, and their available tags

### Example

```ts
import {
  Configuration,
  KbApi,
} from '';
import type { ListKbMemoriesRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new KbApi(config);

  try {
    const data = await api.listKbMemories();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**{ [key: string]: Memory; }**](Memory.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | memories are returned |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## searchKb

> Array&lt;KnowledgeChunk&gt; searchKb(memory, q, tag)



searches knowledge chunks in a memory slot

### Example

```ts
import {
  Configuration,
  KbApi,
} from '';
import type { SearchKbRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new KbApi(config);

  const body = {
    // string
    memory: memory_example,
    // string
    q: q_example,
    // Array<string> (optional)
    tag: ...,
  } satisfies SearchKbRequest;

  try {
    const data = await api.searchKb(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **memory** | `string` |  | [Defaults to `undefined`] |
| **q** | `string` |  | [Defaults to `undefined`] |
| **tag** | `Array<string>` |  | [Optional] |

### Return type

[**Array&lt;KnowledgeChunk&gt;**](KnowledgeChunk.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | knowledge chunks are returned |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## submitDocument

> submitDocument(memory, document, document2)



submits a new document to a memory slot

### Example

```ts
import {
  Configuration,
  KbApi,
} from '';
import type { SubmitDocumentRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new KbApi(config);

  const body = {
    // string
    memory: memory_example,
    // string
    document: document_example,
    // Document (optional)
    document2: ...,
  } satisfies SubmitDocumentRequest;

  try {
    const data = await api.submitDocument(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **memory** | `string` |  | [Defaults to `undefined`] |
| **document** | `string` |  | [Defaults to `undefined`] |
| **document2** | [Document](Document.md) |  | [Optional] |

### Return type

`void` (Empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | knowledge is acquired |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


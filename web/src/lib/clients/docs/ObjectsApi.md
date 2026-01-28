# ObjectsApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createObject**](ObjectsApi.md#createobject) | **POST** /objects/{memory} |  |
| [**deleteObjectByName**](ObjectsApi.md#deleteobjectbyname) | **DELETE** /objects/{memory}/_by-name |  |
| [**getObjectByName**](ObjectsApi.md#getobjectbyname) | **GET** /objects/{memory}/_by-name |  |
| [**listObjects**](ObjectsApi.md#listobjects) | **GET** /objects/{memory} |  |
| [**listObjectsMemories**](ObjectsApi.md#listobjectsmemories) | **GET** /objects/_memories |  |



## createObject

> createObject(memory, dataObject)



### Example

```ts
import {
  Configuration,
  ObjectsApi,
} from '';
import type { CreateObjectRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ObjectsApi(config);

  const body = {
    // string
    memory: memory_example,
    // DataObject (optional)
    dataObject: ...,
  } satisfies CreateObjectRequest;

  try {
    const data = await api.createObject(body);
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
| **dataObject** | [DataObject](DataObject.md) |  | [Optional] |

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
| **201** | object is created |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## deleteObjectByName

> deleteObjectByName(memory, name)



### Example

```ts
import {
  Configuration,
  ObjectsApi,
} from '';
import type { DeleteObjectByNameRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ObjectsApi(config);

  const body = {
    // string
    memory: memory_example,
    // string
    name: name_example,
  } satisfies DeleteObjectByNameRequest;

  try {
    const data = await api.deleteObjectByName(body);
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
| **name** | `string` |  | [Defaults to `undefined`] |

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
| **204** | object is deleted |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getObjectByName

> DataObject getObjectByName(memory, name)



### Example

```ts
import {
  Configuration,
  ObjectsApi,
} from '';
import type { GetObjectByNameRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ObjectsApi(config);

  const body = {
    // string
    memory: memory_example,
    // string
    name: name_example,
  } satisfies GetObjectByNameRequest;

  try {
    const data = await api.getObjectByName(body);
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
| **name** | `string` |  | [Defaults to `undefined`] |

### Return type

[**DataObject**](DataObject.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | object is returned |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listObjects

> Array&lt;DataObject&gt; listObjects(memory)



### Example

```ts
import {
  Configuration,
  ObjectsApi,
} from '';
import type { ListObjectsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ObjectsApi(config);

  const body = {
    // string
    memory: memory_example,
  } satisfies ListObjectsRequest;

  try {
    const data = await api.listObjects(body);
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

[**Array&lt;DataObject&gt;**](DataObject.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | objects are returned |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listObjectsMemories

> Array&lt;string&gt; listObjectsMemories()



### Example

```ts
import {
  Configuration,
  ObjectsApi,
} from '';
import type { ListObjectsMemoriesRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new ObjectsApi(config);

  try {
    const data = await api.listObjectsMemories();
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

**Array<string>**

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


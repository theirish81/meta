# RecipesApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createRecipe**](RecipesApi.md#createrecipe) | **POST** /recipes/{memory} |  |
| [**deleteRecipe**](RecipesApi.md#deleterecipe) | **DELETE** /recipes/{memory}/{recipeId} |  |
| [**listRecipesMemories**](RecipesApi.md#listrecipesmemories) | **GET** /recipes/_memories |  |
| [**searchRecipes**](RecipesApi.md#searchrecipes) | **GET** /recipes/{memory} |  |
| [**updateRecipe**](RecipesApi.md#updaterecipe) | **POST** /recipes/{memory}/{recipeId} |  |



## createRecipe

> Recipe createRecipe(memory, recipeRequest)



creates a new recipe in a memory slot

### Example

```ts
import {
  Configuration,
  RecipesApi,
} from '';
import type { CreateRecipeRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new RecipesApi(config);

  const body = {
    // string
    memory: memory_example,
    // RecipeRequest (optional)
    recipeRequest: ...,
  } satisfies CreateRecipeRequest;

  try {
    const data = await api.createRecipe(body);
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
| **recipeRequest** | [RecipeRequest](RecipeRequest.md) |  | [Optional] |

### Return type

[**Recipe**](Recipe.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | meta is created |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## deleteRecipe

> deleteRecipe(memory, recipeId)



deletes a specific recipe

### Example

```ts
import {
  Configuration,
  RecipesApi,
} from '';
import type { DeleteRecipeRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new RecipesApi(config);

  const body = {
    // string
    memory: memory_example,
    // string
    recipeId: 38400000-8cf0-11bd-b23e-10b96e4ef00d,
  } satisfies DeleteRecipeRequest;

  try {
    const data = await api.deleteRecipe(body);
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
| **recipeId** | `string` |  | [Defaults to `undefined`] |

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
| **204** | meta is deleted |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listRecipesMemories

> { [key: string]: Memory; } listRecipesMemories()



lists all the recipes memory slots, and their available tags

### Example

```ts
import {
  Configuration,
  RecipesApi,
} from '';
import type { ListRecipesMemoriesRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new RecipesApi(config);

  try {
    const data = await api.listRecipesMemories();
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
| **200** | meta tags are returned |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## searchRecipes

> Array&lt;Recipe&gt; searchRecipes(memory, tag, q)



searches recipes in a memory slot

### Example

```ts
import {
  Configuration,
  RecipesApi,
} from '';
import type { SearchRecipesRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new RecipesApi(config);

  const body = {
    // string
    memory: memory_example,
    // Array<string> (optional)
    tag: ...,
    // string (optional)
    q: q_example,
  } satisfies SearchRecipesRequest;

  try {
    const data = await api.searchRecipes(body);
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
| **tag** | `Array<string>` |  | [Optional] |
| **q** | `string` |  | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;Recipe&gt;**](Recipe.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | meta are returned |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## updateRecipe

> Recipe updateRecipe(memory, recipeId, recipeRequest)



updates a specific recipe

### Example

```ts
import {
  Configuration,
  RecipesApi,
} from '';
import type { UpdateRecipeRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new RecipesApi(config);

  const body = {
    // string
    memory: memory_example,
    // string
    recipeId: 38400000-8cf0-11bd-b23e-10b96e4ef00d,
    // RecipeRequest (optional)
    recipeRequest: ...,
  } satisfies UpdateRecipeRequest;

  try {
    const data = await api.updateRecipe(body);
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
| **recipeId** | `string` |  | [Defaults to `undefined`] |
| **recipeRequest** | [RecipeRequest](RecipeRequest.md) |  | [Optional] |

### Return type

[**Recipe**](Recipe.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | meta is updated |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


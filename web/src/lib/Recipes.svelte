
<!--
  - Copyright (C) 2026 Simone Pezzano
  -
  - This program is free software: you can redistribute it and/or modify
  - it under the terms of the GNU Affero General Public License as
  - published by the Free Software Foundation, either version 3 of the
  - License, or (at your option) any later version.
  -
  - This program is distributed in the hope that it will be useful,
  - but WITHOUT ANY WARRANTY; without even the implied warranty of
  - MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  - GNU Affero General Public License for more details.
  -
  - You should have received a copy of the GNU Affero General Public License
  - along with this program.  If not, see <https://www.gnu.org/licenses/>.
  -->


<script lang="ts">
    import {
        Label,
        Select,
        MultiSelect,
        Input,
        Textarea,
        Button,
        Modal,
        type SelectOptionType,
        Alert
    } from "flowbite-svelte";
    import {Card, Badge, Tags} from 'flowbite-svelte';
    import { TrashBinSolid } from "flowbite-svelte-icons";

    import type {Memory, Recipe} from "./clients";
    import {getRecipesClient} from "./session.svelte.js";
    import {onMount} from "svelte";

    const client = getRecipesClient()
    let selectedMemory = $state('')
    let memories: {[key: string]: Memory} = $state({})
    let memoriesList = $derived.by(() => {
        return Object.keys(memories).map(key => ({"name": key, "value": key}))
    })
    let tagsSelected = $state([])
    let query = $state('')
    let items: Recipe[] = $state([])
    let error = $state('')
    let formOpen = $state(false)
    let formMemory = $state('')
    let selectedItem: Recipe = $state(<Recipe>{})
    onMount(() => {
        listMemories()
    })

    function getMemoryTagsSelect(): SelectOptionType<any>[]{
        if (!memories || !selectedMemory) {
            return []
        }
        return memories[selectedMemory]?.availableTags.map(key => ({"name": key, "value": key}))
    }

    async function listMemories(){
        try {
            memories = await client.listRecipesMemories()
        }catch(e: any) {
            error = "could not list memories: "+e.message
        }
    }
    async function searchRecipes(){
        const tags = tagsSelected.length > 0 ? tagsSelected : undefined
        const q = query.length > 0 ? query : undefined
        try {
            items = await client.searchRecipes({
                memory: selectedMemory,
                tag: tags,
                q: q
            })
        }catch(e: any) {
            error = "could not search recipes: "+e.message
        }
    }

    function isSearchSubmitDisabled(): boolean {
        return !selectedMemory
    }

    function edit(recipe: Recipe) {
        formMemory = selectedMemory
        selectedItem = {...recipe}
        formOpen = true
    }
    function create() {
        formMemory = selectedMemory
        selectedItem = <Recipe>{}
        selectedItem.tags = []
        formOpen = true
    }

    async function save() {
        try {
            if (selectedItem.id) {
                await getRecipesClient().updateRecipe({
                    recipeId: selectedItem.id,
                    memory: formMemory,
                    recipeRequest: selectedItem
                })
            } else {
                await getRecipesClient().createRecipe({
                    memory: formMemory,
                    recipeRequest: selectedItem
                })
            }
        }catch(e: any) {
            error = "could not save recipe: "+e.message
        }
        formOpen = false
        await listMemories()
        selectedMemory = formMemory
        await searchRecipes()
    }

    async function deleteRecipe() {
        try {

            await getRecipesClient().deleteRecipe({
                recipeId: selectedItem.id,
                memory: formMemory,
            })
        }catch(e: any) {
            error = "could not delete recipe: "+e.message
        }
        formOpen = false
        await listMemories()
        if (Object.keys(memories).includes(formMemory)) {
            selectedMemory = formMemory
            await searchRecipes()
        } else {
            items = []
        }
    }

    function onModalClose() {
        selectedItem = <Recipe>{}
        formOpen = false
        formMemory = ''
    }

    function isFormSaveDisabled(): boolean {
        return !(formMemory && selectedItem.name && selectedItem.description && selectedItem.content && selectedItem.tags)
    }
</script>


<Modal open={formOpen} onclose={onModalClose} class="w-full max-w-5xl" permanent>

    <div class="mb-6">
        <Label for="name" class="mb-2 block">Memory</Label>
        <Input id="name" size="lg" disabled={selectedItem.id != null } bind:value={formMemory} />
    </div>
    <div class="mb-6">
        <Label for="name" class="mb-2 block">Name</Label>
        <Input id="name" size="lg"  bind:value={selectedItem.name} />
    </div>
    <div class="mb-6">
        <Label for="tags" class="mb-2 block">Tags</Label>
        <Tags bind:value={selectedItem.tags} unique={true} />
    </div>
    <div class="mb-6">
        <Label for="description" class="mb-2 block">Description</Label>
        <Input id="description" size="lg"  bind:value={selectedItem.description} />
    </div>
    <div class="mb-6">
        <Label for="content" class="mb-2 block">Content</Label>
        <Textarea id="content" class="w-full h-80" bind:value={selectedItem.content} />
    </div>
    <div class="flex justify-between">
        <div>
            <Button class="mt-4" onclick={() => save()} disabled={isFormSaveDisabled()}>Save</Button>
            <Button class="mt-4" onclick={onModalClose}>Cancel</Button>
        </div>
        {#if selectedItem.id}
            <Button class="mt-4" onclick={deleteRecipe}>
                <TrashBinSolid />
            </Button>
        {/if}

    </div>

</Modal>

<div>
    <Alert color="red" alertStatus={error.length > 0}>{error}</Alert>
</div>
<div class="flex flex-wrap gap-4 items-end">

    <div class="flex-1">
        <Label for="memory">Memory *</Label>
        <Select name="memory" class="mt-2" items={memoriesList} bind:value={selectedMemory} onchange={() => tagsSelected = []}/>
    </div>
    <div class="flex-1">
        <Label for="tags">Tags</Label>
        <MultiSelect name="tags" items={getMemoryTagsSelect()} bind:value={tagsSelected} />
    </div>
    <div class="flex-1">
        <Label for="query">Relevance Query</Label>
        <Input name="query" type="text" bind:value={query}/>
    </div>
    <div class="flex-1">
        <Button onclick={searchRecipes} disabled={isSearchSubmitDisabled()}>Search</Button>
    </div>
</div>
<div class="w-full p-20 flex flex-wrap gap-4">
    <Card class="p-4 sm:p-6 md:p-8 h-50" onclick={() => create()}>
        <div class="flex items-center justify-center text-9xl text-gray-700 dark:text-gray-400 h-full">
           +
       </div>
    </Card>
    {#each items as recipe}
        <Card class="p-4 sm:p-6 md:p-8 h-50" onclick={() => edit(recipe)}>
            <h5 class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">{recipe.name}</h5>
            <p class="leading-tight font-normal text-gray-700 dark:text-gray-400">{recipe.description}</p>
            <p>
                {#each recipe.tags as tag}
                    <Badge>{tag}</Badge>
                {/each}
            </p>
        </Card>
    {/each}
</div>
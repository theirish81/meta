
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
        Input,
        Textarea,
        Button,
        Modal,
        Alert
    } from "flowbite-svelte";
    import {Card} from 'flowbite-svelte';
    import { TrashBinSolid } from "flowbite-svelte-icons";

    import type {DataObject } from "./clients";
    import {getObjectsClient } from "./session.svelte.js";
    import {onMount} from "svelte";

    const client = getObjectsClient()
    let selectedMemory = $state('')
    let memories: string[] = $state([])
    let objects: DataObject[] = $state([])
    let error = $state('')
    let formOpen = $state(false)
    let formMemory = $state('')
    let selectedItem: DataObject = $state(<DataObject>{})

    let memoriesList = $derived.by(() => {
        return memories.map(key => ({"name": key, "value": key}))
    })
    onMount(() => {
        listMemories()
    })


    async function listMemories(){
        try {
            memories = await client.listObjectsMemories()
        }catch(e: any) {
            error = "could not list memories: "+e.message
        }
    }


    async function listObjects() {
        try{
            objects = await client.listObjects({memory: selectedMemory})
        }catch(e: any){
            error = "could not list objects: "+e.message
        }
    }


    function edit(object: DataObject) {
        formMemory = selectedMemory
        selectedItem = {...object}
        formOpen = true
    }

    async function deleteObject() {
        try {

            await getObjectsClient().deleteObjectByName({
                name: selectedItem.name,
                memory: formMemory,
            })
        }catch(e: any) {
            error = "could not delete object: "+e.message
        }
        formOpen = false
        await listMemories()
        if (Object.keys(memories).includes(formMemory)) {
            selectedMemory = formMemory
            await listObjects()
        } else {
            objects = []
        }
    }

    function onModalClose() {
        selectedItem = <DataObject>{}
        formOpen = false
        formMemory = ''
    }
</script>


<Modal open={formOpen} onclose={onModalClose} class="w-full max-w-5xl" permanent>
    <div class="mb-6">
        <Label for="name" class="mb-2 block">Name</Label>
        <Input id="name" size="lg"  bind:value={selectedItem.name} readonly />
    </div>
    <div class="mb-6">
        <Label for="contentType" class="mb-2 block">Content-Type</Label>
        <Input id="contentType" size="lg"  bind:value={selectedItem.contentType} readonly />
    </div>
    <div class="mb-6">
        <Label for="content" class="mb-2 block">Content</Label>
        <Textarea id="content" class="w-full h-80" bind:value={selectedItem.content} readonly/>
    </div>
    <div class="flex justify-between">
        <div>
            <Button class="mt-4" onclick={onModalClose}>Cancel</Button>
        </div>
        <Button class="mt-4" onclick={deleteObject}>
            <TrashBinSolid />
        </Button>
    </div>

</Modal>

<div>
    <Alert color="red" alertStatus={error.length > 0}>{error}</Alert>
</div>
<div class="flex flex-wrap gap-4 items-end">
    <div class="flex-1">
        <Label for="memory">Memory *</Label>
        <Select name="memory" class="mt-2" items={memoriesList} bind:value={selectedMemory} onchange={() => listObjects()}/>
    </div>
</div>
<div class="w-full p-20 flex flex-wrap gap-4">
    {#each objects as obj}
        <Card class="p-4 sm:p-6 md:p-8 h-50" onclick={() => edit(obj)}>
            <h5 class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">{obj.name}</h5>
            <p class="leading-tight font-normal text-gray-700 dark:text-gray-400 line-clamp-4">{obj.content}</p>
        </Card>
    {/each}
</div>
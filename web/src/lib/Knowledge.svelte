
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
        Button,
        Modal,
        type SelectOptionType,
        Textarea,
        Tags, Alert
    } from "flowbite-svelte";
    import { Badge, Table, TableHead, TableHeadCell, TableBody, TableBodyRow, TableBodyCell} from 'flowbite-svelte';


    import type {KnowledgeChunk, Memory, Document} from "./clients";
    import {getKnowledgeClient} from "./session.svelte.js";
    import {onMount} from "svelte";
    import {TrashBinSolid} from "flowbite-svelte-icons";

    const client = getKnowledgeClient()
    let error = $state('')
    let selectedMemory = $state('')
    let memories: {[key: string]: Memory} = $state({})
    let memoriesList = $derived.by(() => {
        return Object.keys(memories).map(key => ({"name": key, "value": key}))
    })
    let tagsSelected = $state([])
    let query = $state('')
    let items: KnowledgeChunk[] = $state([])

    let formMemory = $state('')
    let formDocumentData: Document = $state(<Document>{})
    let formDocumentName = $state('')
    let formOpen = $state(false)

    let panelDocuments: string[] = $state([])
    let panelDocumentsOpen = $state(false)


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
            memories = await client.listKbMemories()
        }catch(e: any) {
            error = "could not list memories: "+e.message
        }
    }

    async function searchItems(){
        const tags = tagsSelected.length > 0 ? tagsSelected : undefined
        try {
            items = await client.searchKb({
                memory: selectedMemory,
                tag: tags,
                q: query
            })
        }catch(e: any) {
            error = "could not search items: "+e.message
        }
    }

    async function listDocuments() {
        try {
            panelDocuments = await client.listDocuments({
                memory: selectedMemory,
            })
        }catch(e: any) {
            error = "could not list documents: "+e.message
        }
    }

    function isSearchSubmitDisabled(): boolean {
        return !selectedMemory || query.length == 0
    }

    async function showDocuments() {
        await listDocuments()
        panelDocumentsOpen = true
    }
    async function deleteDocument(document: string) {
        try {
            await client.deleteDocument({
                memory: selectedMemory,
                document: document
            })
        }catch(e: any) {
            error = "could not delete document: "+e.message
        }
        await listMemories()
        selectedMemory = ""
        items = []
        panelDocumentsOpen = false
    }

    async function save() {
        try {
            await client.submitDocument({
                memory: formMemory,
                document: formDocumentName,
                document2: formDocumentData
            })
        }catch(e: any) {
            error = "could not save document: "+e.message
        }
        await listMemories()
        formOpen = false
    }

    function isFormSaveDisabled() {
        return !formDocumentName || !formDocumentData.content || !formMemory || formDocumentData.tags.length == 0
    }

    function openForm() {
        formDocumentName = ""
        formMemory = selectedMemory
        formDocumentData = <Document>{}
        formDocumentData.tags = []
        formOpen = true

    }

</script>


<Modal class="w-full max-w-5xl p-10" open={panelDocumentsOpen}>
        {#each panelDocuments as document}

            <div class="flex justify-between items-end">
                <div class="mt-4 text-xl">
                    {document}
                </div>
                <div>
                    <Button class="mt-4" onclick={() => deleteDocument(document)}>
                        <TrashBinSolid />
                    </Button>
                </div>
            </div>
        {/each}
</Modal>

<Modal class="w-full max-w-5xl p-10" open={formOpen}>
    <div class="mb-6">
        <Label for="name" class="mb-2 block">Memory</Label>
        <Input id="name" size="lg" bind:value={formMemory}/>
    </div>
    <div class="mb-6">
        <Label for="document" class="mb-2 block">Name</Label>
        <Input id="document" size="lg"  bind:value={formDocumentName} />
    </div>
    <div class="mb-6">
        <Label for="tags" class="mb-2 block">Tags</Label>
        <Tags bind:value={formDocumentData.tags}></Tags>
    </div>
    <div class="mb-6">
        <Label for="content" class="mb-2 block">Content</Label>
        <Textarea id="content" class="w-full h-80" bind:value={formDocumentData.content} />
    </div>

    <div class="flex justify-between">
        <Button class="mt-4" onclick={() => save()} disabled={isFormSaveDisabled()}>Save</Button>
    </div>
</Modal>

<div>
    <Alert color="red" alertStatus={error.length > 0}>{error}</Alert>
</div>
<div class="flex flex-wrap gap-4 items-end">
    <div class="flex-1">
        <Label for="memory">Memories</Label>
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
        <Button onclick={searchItems} disabled={isSearchSubmitDisabled()}>Search</Button>
        <Button onclick={openForm}>+</Button>
    </div>
</div>
<div class="w-full mt-4">
    <Table class="table-fixed w-full">
        <TableHead>
            <TableHeadCell class="w-1/6">Document</TableHeadCell>
            <TableHeadCell class="w-1/6">Tags</TableHeadCell>
            <TableHeadCell class="w-4/6">Chunk</TableHeadCell>
        </TableHead>
        <TableBody>
            {#each items as item}
                <TableBodyRow>
                    <TableBodyCell>
                        <a href="#" onclick={showDocuments}>{item.document}</a>
                    </TableBodyCell>
                    <TableBodyCell>
                        {#each item.tags as t}
                            <Badge>{t}</Badge>
                        {/each}
                    </TableBodyCell>
                    <TableBodyCell class="whitespace-normal break-words break-all">{item.chunk}</TableBodyCell>
                </TableBodyRow>
            {/each}
        </TableBody>
    </Table>
</div>
<script lang="ts">
  import { onMount } from 'svelte';
  import { Tabs, TabItem, Button, Input } from 'flowbite-svelte';
  import './app.css'
  import Recipes from "./lib/Recipes.svelte";
  import Knowledge from "./lib/Knowledge.svelte";

  let authorizationKey: string | null = null;
  let isAuthenticated = $state(false);
  let tokenInput = $state('');

  let kbMemories = $state({})
  let multiSelected: string[] = [];

  onMount(() => {
    authorizationKey = localStorage.getItem('authorizationKey');
    if (authorizationKey) {
      isAuthenticated = true;
    }

  });

  function handleLogin(e: Event) {
    e.preventDefault();
    if (tokenInput) {
      localStorage.setItem('authorizationKey', tokenInput);
      authorizationKey = tokenInput;
      isAuthenticated = true;
    }
  }

  function logout() {
      localStorage.removeItem('authorizationKey');
      authorizationKey = ""
      isAuthenticated = false;
  }

</script>

<main class="p-4">
  {#if isAuthenticated}
      <div class="flex justify-between">
          <h1 class="text-4xl font-bold">Meta</h1>
          <Button onclick={logout}>Log out</Button>
      </div>
    <Tabs>
      <TabItem title="Recipes">
          <Recipes/>
      </TabItem>
      <TabItem title="Knowledge">
        <Knowledge/>
      </TabItem>
    </Tabs>
  {:else}
    <div class="flex items-center justify-center h-screen">
      <div class="w-full max-w-xs">
        <form
          class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4"
          onsubmit={handleLogin}
        >
          <div class="mb-4">
            <label
              class="block text-gray-700 text-sm font-bold mb-2"
              for="token"
            >
              Authorization Token
            </label>
            <Input
              id="token"
              type="password"
              bind:value={tokenInput}
            />
          </div>
          <div class="flex items-center justify-between">
            <Button type="submit">Sign In</Button>
          </div>
        </form>
      </div>
    </div>
  {/if}
</main>
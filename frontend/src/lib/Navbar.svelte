<script>
  import { onMount } from 'svelte';
  import { link } from 'svelte-routing';
  import { currentUser, pb } from './pocketbase';

  let apiStatus = 'pending'; // 'pending', 'online', 'offline'
  let statusText = 'Checking API status...';

  async function checkApiHealth() {
    try {
      const response = await fetch('/api/health');
      if (response.ok) {
        const data = await response.json();
        if (data.status === 'ok') {
          apiStatus = 'online';
          statusText = 'API is online';
        } else {
          apiStatus = 'offline';
          statusText = 'API returned an unexpected status';
        }
      } else {
        apiStatus = 'offline';
        statusText = 'API is offline or unreachable';
      }
    } catch (error) {
      apiStatus = 'offline';
      statusText = 'Failed to connect to the API';
    }
  }

  onMount(() => {
    checkApiHealth();
    // Check health every 30 seconds
    const interval = setInterval(checkApiHealth, 30000); 
    return () => clearInterval(interval);
  });

  function logout() {
    pb.authStore.clear();
  }
</script>

<nav class="bg-white dark:bg-gray-800 shadow-md">
  <div class="container mx-auto px-6 py-3">
    <div class="flex items-center justify-between">
      <div class="flex items-center">
        <a href="/" use:link class="flex items-center text-xl font-semibold text-gray-700 dark:text-white">
          <img src="/vite.svg" alt="Logo" class="w-8 h-8 mr-2" />
          <span>Receipt Scanner</span>
        </a>
      </div>

      <div class="flex items-center">
        {#if apiStatus === 'offline'}
          <div class="relative group mr-6">
            <div class="flex items-center cursor-pointer">
              <span class="mr-2 h-3 w-3 rounded-full bg-red-500"></span>
              <span class="text-sm text-red-500 dark:text-red-400">API Offline</span>
            </div>
            <div class="absolute bottom-full mb-2 w-48 hidden group-hover:block bg-gray-700 text-white text-xs rounded py-1 px-2 text-center z-10">
              {statusText}
            </div>
          </div>
        {/if}

        {#if $currentUser}
          <span class="text-gray-800 dark:text-gray-200 mx-3">Welcome, {$currentUser.name || $currentUser.email}</span>
          <a href="/settings" use:link class="text-gray-800 dark:text-gray-200 hover:text-indigo-600 dark:hover:text-indigo-400 mx-3">Settings</a>
          <button on:click={logout} class="text-gray-800 dark:text-gray-200 hover:text-indigo-600 dark:hover:text-indigo-400 mx-3">Logout</button>
        {:else}
          <a href="/login" use:link class="text-gray-800 dark:text-gray-200 hover:text-indigo-600 dark:hover:text-indigo-400 mx-3">Login</a>
          <a href="/register" use:link class="text-gray-800 dark:text-gray-200 hover:text-indigo-600 dark:hover:text-indigo-400 mx-3">Register</a>
        {/if}
      </div>
    </div>
  </div>
</nav> 
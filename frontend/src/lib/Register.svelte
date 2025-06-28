<script>
  import { onMount } from 'svelte';

  let name = '';
  let email = '';
  let password = '';
  let passwordConfirm = '';
  let message = '';
  let messageType = ''; // 'success' or 'error'

  async function handleRegister() {
    if (password !== passwordConfirm) {
      message = "Passwords do not match.";
      messageType = 'error';
      return;
    }

    message = 'Creating account...';
    messageType = 'info';

    try {
      const response = await fetch('/api/users/create', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, email, password, passwordConfirm }),
      });

      const result = await response.json();

      if (response.ok) {
        message = 'Account created successfully! You can now log in.';
        messageType = 'success';
        name = '';
        email = '';
        password = '';
        passwordConfirm = '';
      } else {
        message = `Registration failed: ${result.message || 'Unknown error'}`;
        messageType = 'error';
      }
    } catch (error) {
      message = `An error occurred: ${error.message}`;
      messageType = 'error';
    }
  }
</script>

<div class="flex items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-900">
  <div class="w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md dark:bg-gray-800">
    <h1 class="text-3xl font-bold text-center text-gray-900 dark:text-white">Create an Account</h1>
    
    <form on:submit|preventDefault={handleRegister} class="space-y-6">
      <div>
        <label for="name" class="block mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">Name</label>
        <input type="text" id="name" bind:value={name} required class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white">
      </div>
      <div>
        <label for="email" class="block mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">Email Address</label>
        <input type="email" id="email" bind:value={email} required class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white">
      </div>
      <div>
        <label for="password" class="block mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">Password</label>
        <input type="password" id="password" bind:value={password} required class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white">
      </div>
      <div>
        <label for="passwordConfirm" class="block mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">Confirm Password</label>
        <input type="password" id="passwordConfirm" bind:value={passwordConfirm} required class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white">
      </div>
      <button type="submit" class="w-full px-4 py-2 text-lg font-semibold text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-300">
        Register
      </button>
    </form>

    {#if message}
      <div class="text-center p-4 rounded-lg"
        class:bg-green-100={messageType === 'success'} class:text-green-800={messageType === 'success'}
        class:bg-red-100={messageType === 'error'} class:text-red-800={messageType === 'error'}
        class:bg-blue-100={messageType === 'info'} class:text-blue-800={messageType === 'info'}
      >
        {message}
      </div>
    {/if}

    <p class="text-sm text-center text-gray-600 dark:text-gray-400">
      Already have an account? <a href="/login" class="font-medium text-indigo-600 hover:underline dark:text-indigo-400">Log in</a>
    </p>
  </div>
</div> 
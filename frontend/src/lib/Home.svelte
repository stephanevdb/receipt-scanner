<script>
  import { currentUser } from './pocketbase';
  import { link } from 'svelte-routing';

  let selectedFile = null;
  let message = '';
  let receiptUrl = ''; // To store the URL of the uploaded image

  function handleFileSelect(event) {
    selectedFile = event.target.files[0];
    message = '';
    receiptUrl = '';
    if (selectedFile) {
      receiptUrl = URL.createObjectURL(selectedFile);
    }
  }

  async function handleUpload() {
    if (!selectedFile) {
      message = 'Please select a file first.';
      return;
    }

    const formData = new FormData();
    formData.append('receipt', selectedFile);

    message = 'Uploading...';

    try {
      const response = await fetch('/api/receipts/upload', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: 'An unknown error occurred.' }));
        message = `Upload failed: ${errorData.message}`;
        return;
      }

      const result = await response.json();

      if (response.ok) {
        message = `Upload successful: ${result.filename}`;
      } else {
        message = `Upload failed: ${result.message}`;
      }
    } catch (error) {
      message = `An error occurred: ${error.message}`;
    }
  }
</script>

<main class="text-gray-800 dark:text-gray-200 transition-colors duration-300">
  <div class="container mx-auto px-4 py-12">
    <header class="text-center mb-12">
      <h1 class="text-5xl font-extrabold text-gray-900 dark:text-white">Receipt Scanner</h1>
      <p class="mt-4 text-lg text-gray-600 dark:text-gray-400">
        {#if $currentUser}
          Upload a receipt to get started.
        {:else}
          Please log in or register to upload and manage your receipts.
        {/if}
      </p>
    </header>

    {#if $currentUser}
      <div class="w-full max-w-2xl mx-auto">
        <div
          class="bg-white dark:bg-gray-800 rounded-xl shadow-lg transition-shadow duration-300 ease-in-out hover:shadow-2xl"
        >
          <div class="p-8">
            <label
              for="file-upload"
              class="relative flex flex-col items-center justify-center w-full h-64 border-2 border-gray-300 dark:border-gray-600 border-dashed rounded-lg cursor-pointer bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 transition-all duration-300"
            >
              {#if receiptUrl}
                <img src={receiptUrl} alt="Selected receipt" class="h-full w-full object-contain rounded-lg p-2" />
              {:else}
                <div class="flex flex-col items-center justify-center text-center">
                  <svg class="w-10 h-10 mb-4 text-gray-400 dark:text-gray-500" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 16"><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 13h3a3 3 0 0 0 0-6h-.025A5.56 5.56 0 0 0 16 6.5 5.5 5.5 0 0 0 5.207 5.021C5.137 5.017 5.071 5 5 5a4 4 0 0 0 0 8h2.167M10 15V6m0 0L8 8m2-2 2 2" /></svg>
                  <p class="mb-2 text-sm text-gray-500 dark:text-gray-400">
                    <span class="font-semibold">Click to upload</span> or drag and drop
                  </p>
                  <p class="text-xs text-gray-500 dark:text-gray-400">PNG or JPG</p>
                </div>
              {/if}
              <input id="file-upload" type="file" class="hidden" on:change={handleFileSelect} accept="image/png, image/jpeg, image/jpg" />
            </label>

            {#if selectedFile}
              <div class="mt-6 text-center">
                <p class="text-gray-600 dark:text-gray-300">Selected file: <span class="font-medium text-gray-800 dark:text-white">{selectedFile.name}</span></p>
              </div>
            {/if}

            <div class="mt-6">
              <button
                on:click={handleUpload}
                class="w-full px-6 py-3 text-lg font-semibold text-white bg-indigo-600 rounded-lg shadow-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:bg-indigo-400 disabled:cursor-not-allowed transition-all duration-300 ease-in-out"
                disabled={!selectedFile}
              >
                Upload
              </button>
            </div>

            {#if message}
              <p class="mt-6 text-center text-sm"
                class:text-green-600={message.startsWith('Upload successful')}
                class:dark:text-green-400={message.startsWith('Upload successful')}
                class:text-red-600={!message.startsWith('Upload successful')}
                class:dark:text-red-400={!message.startsWith('Upload successful')}
              >{message}</p>
            {/if}
          </div>
        </div>
      </div>
    {:else}
      <div class="text-center">
        <p class="text-xl text-gray-700 dark:text-gray-300">Welcome to Receipt Scanner!</p>
        <div class="mt-6">
          <a href="/login" use:link class="px-6 py-3 text-lg font-semibold text-white bg-indigo-600 rounded-lg shadow-md hover:bg-indigo-700 transition-all duration-300 ease-in-out">
            Log In
          </a>
          <a href="/register" use:link class="ml-4 px-6 py-3 text-lg font-semibold text-indigo-600 bg-white border border-indigo-600 rounded-lg shadow-md hover:bg-indigo-50 transition-all duration-300 ease-in-out">
            Register
          </a>
        </div>
      </div>
    {/if}
		
		  <footer class="text-center mt-12 text-sm text-gray-500 dark:text-gray-400">
			  <p>&copy; {new Date().getFullYear()} Receipt Scanner. All Rights Reserved.</p>
		  </footer>
  </div>
</main> 
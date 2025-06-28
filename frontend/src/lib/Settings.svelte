<script>
    import { onMount } from 'svelte';
    import { pb } from './pocketbase';

    let user;
    let iban = '';
    let bic = '';

    let newGroupName = '';
    let joinGroupCode = '';
    
    let currentGroup = null;

    async function fetchCurrentGroup(groupId) {
        if (!groupId) {
            currentGroup = null;
            return;
        }
        try {
            currentGroup = await pb.collection('friend_groups').getOne(groupId);
        } catch (error) {
            console.error('Failed to fetch current group:', error);
            alert('Could not load your group details.');
            // The group might have been deleted, so we should clear it from the user
            if (error.status === 404) {
                 try {
                    await pb.collection('users').update(user.id, { 'friend_group': null });
                    user.friend_group = null;
                 } catch (updateError) {
                    console.error('Failed to clear non-existent group from user:', updateError);
                 }
            }
        }
    }

    onMount(async () => {
        user = pb.authStore.model;
        if (user) {
            iban = user.iban || '';
            bic = user.bic || '';
            if (user.friend_group) {
                await fetchCurrentGroup(user.friend_group);
            }
        }
    });

    async function updateSettings() {
        if (!user) return;
        try {
            await pb.collection('users').update(user.id, {
                iban,
                bic,
            });
            alert('Settings updated successfully!');
        } catch (error) {
            console.error('Failed to update settings:', error);
            alert('Failed to update settings.');
        }
    }

    async function createGroup() {
        if (!newGroupName.trim()) {
            alert('Please enter a group name.');
            return;
        }
        try {
            const newGroup = await pb.collection('friend_groups').create({
                name: newGroupName,
            });
            
            await joinGroup(newGroup.id);

            newGroupName = '';
            alert(`Group '${newGroup.name}' created and joined successfully!`);
        } catch (error) {
            console.error('Failed to create group:', error);
            alert('Failed to create group.');
        }
    }

    async function joinGroup(groupId) {
        const idToJoin = groupId || joinGroupCode;
        if (!idToJoin.trim()) {
            alert('Please enter a group code to join.');
            return;
        }
        if (!user) return;

        try {
            // Check if group exists before trying to join
            await pb.collection('friend_groups').getOne(idToJoin);

            await pb.collection('users').update(user.id, {
                'friend_group': idToJoin
            });

            // Refresh the auth store to get the latest user details
            await pb.collection('users').authRefresh();
            user = pb.authStore.model; // update local user model

            await fetchCurrentGroup(idToJoin);
            joinGroupCode = '';
            alert('Successfully joined group!');
            
        } catch (err) {
            console.error('Failed to join group:', err);
            if (err.status === 404) {
                alert('Group not found. Please check the code and try again.');
            } else {
                alert('Error joining group.');
            }
        }
    }

    async function leaveGroup() {
        if (!user || !currentGroup) return;

        if (!confirm(`Are you sure you want to leave the group "${currentGroup.name}"?`)) {
            return;
        }

        try {
            await pb.collection('users').update(user.id, {
                'friend_group': null
            });
             // Refresh the auth store to get the latest user details
            await pb.collection('users').authRefresh();
            user = pb.authStore.model; // update local user model

            currentGroup = null;
            alert('You have left the group.');
        } catch (err) {
            console.error('Failed to leave group:', err);
            alert('Error leaving group.');
        }
    }
</script>

<div class="container mx-auto p-8">
    <h1 class="text-2xl font-bold mb-6">Settings</h1>

    <div class="bg-white p-6 rounded-lg shadow-md mb-8">
        <h2 class="text-xl font-semibold mb-4">Bank Details</h2>
        <form on:submit|preventDefault={updateSettings}>
            <div class="mb-4">
                <label for="iban" class="block text-gray-700 mb-2">IBAN</label>
                <input type="text" id="iban" bind:value={iban} class="w-full p-2 border border-gray-300 rounded" placeholder="DE89 3704 0044 0532 0130 00">
            </div>
            <div class="mb-4">
                <label for="bic" class="block text-gray-700 mb-2">BIC</label>
                <input type="text" id="bic" bind:value={bic} class="w-full p-2 border border-gray-300 rounded" placeholder="COBADEFFXXX">
            </div>
            <button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">Save Bank Details</button>
        </form>
    </div>

    <div class="bg-white p-6 rounded-lg shadow-md">
        <h2 class="text-xl font-semibold mb-4">Friend Group</h2>

        {#if currentGroup}
            <div class="mb-6 p-4 bg-gray-100 rounded">
                <h3 class="text-lg font-medium mb-2">Your Current Group</h3>
                <p class="text-gray-800"><strong>Name:</strong> {currentGroup.name}</p>
                <p class="text-gray-800 mt-1"><strong>Share Code:</strong> <code class="bg-gray-200 px-2 py-1 rounded text-sm">{currentGroup.id}</code></p>
                <p class="text-xs text-gray-500 mt-2">Share this code with your friends so they can join your group.</p>
                <button on:click={leaveGroup} class="mt-4 bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600 text-sm">Leave Group</button>
            </div>
        {:else}
            <div class="mb-6">
                 <h3 class="text-lg font-medium mb-2">Join a Group with a Code</h3>
                 <div class="flex items-center gap-4">
                    <input type="text" bind:value={joinGroupCode} class="w-full p-2 border border-gray-300 rounded" placeholder="Enter join code">
                    <button on:click={() => joinGroup()} class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600 whitespace-nowrap">Join</button>
                 </div>
            </div>
    
            <div class="border-t pt-6 mt-6">
                <h3 class="text-lg font-medium mb-2">Create a New Group</h3>
                <p class="text-sm text-gray-600 mb-4">No group to join? Create one and invite your friends!</p>
                <div class="flex items-center gap-4">
                    <input type="text" bind:value={newGroupName} class="w-full p-2 border border-gray-300 rounded" placeholder="New group name">
                    <button on:click={createGroup} class="bg-purple-500 text-white px-4 py-2 rounded hover:bg-purple-600 whitespace-nowrap">Create & Join</button>
                </div>
            </div>
        {/if}
    </div>
</div> 
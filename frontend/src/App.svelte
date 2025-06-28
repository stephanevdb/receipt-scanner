<script>
  import { Router, Route } from "svelte-routing";
  import Navbar from './lib/Navbar.svelte';
  import Home from './lib/Home.svelte';
  import Register from './lib/Register.svelte';
  import Login from './lib/Login.svelte';
  import Settings from './lib/Settings.svelte';
</script>

<div class="min-h-screen bg-gray-50 dark:bg-gray-900">
  <Navbar />

  <Router>
    <Route path="/">
      <Home />
    </Route>
    <Route path="/register">
      <Register />
    </Route>
    <Route path="/login">
      <Login />
    </Route>
    <Route path="/settings">
      <Settings />
    </Route>
  </Router>
</div>

<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-xl font-semibold text-white">Manage Posts</h2>
      <button @click="openModal()" class="px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded-lg text-sm font-medium transition-colors flex items-center shadow-sm shadow-primary/20">
        <PlusIcon class="w-4 h-4 mr-2" />
        New Post
      </button>
    </div>

    <div class="bg-dark-surface border border-dark-border rounded-xl shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm text-dark-muted">
          <thead class="text-xs text-white uppercase bg-dark-bg border-b border-dark-border">
            <tr>
              <th class="px-6 py-4 font-medium">Title</th>
              <th class="px-6 py-4 font-medium">Summary</th>
              <th class="px-6 py-4 font-medium">Status</th>
              <th class="px-6 py-4 font-medium">Scheduled For</th>
              <th class="px-6 py-4 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="store.posts.length === 0">
              <td colspan="5" class="px-6 py-8 text-center">No posts found.</td>
            </tr>
            <tr v-for="post in store.posts" :key="post.id" class="border-b border-dark-border hover:bg-dark-bg/50 transition-colors">
              <td class="px-6 py-4 font-medium text-white">{{ post.title }}</td>
              <td class="px-6 py-4">{{ post.summary }}</td>
              <td class="px-6 py-4">
                <StatusBadge :status="post.status" />
              </td>
              <td class="px-6 py-4">{{ new Date(post.post_at).toLocaleString() }}</td>
              <td class="px-6 py-4 text-right">
                <button class="text-primary hover:text-white mr-3 transition-colors">Edit</button>
                <button class="text-red-400 hover:text-red-300 transition-colors">Delete</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useMainStore } from '../stores/main'
import StatusBadge from '../components/StatusBadge.vue'
import { PlusIcon } from '@heroicons/vue/20/solid'

const store = useMainStore()

onMounted(async () => {
  await store.fetchPosts()
})

const openModal = () => {
  alert('Modal implementation depends on requirements. Skeleton ready.')
}
</script>

<template>
  <div class="relative" ref="container">
    <div
      @click="isOpen = !isOpen"
      class="w-full text-sm bg-light-bg border border-light-border rounded-lg p-2.5 cursor-pointer flex justify-between items-center transition-colors"
      :class="[isOpen ? 'ring-2 ring-primary border-primary' : 'hover:border-primary']"
    >
      <span :class="{'text-light-muted': !selectedOption}" class="truncate">
        {{ selectedOption ? selectedOption.label : placeholder }}
      </span>
      <svg class="w-4 h-4 text-light-muted transition-transform" :class="{'rotate-180': isOpen}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
      </svg>
    </div>

    <!-- Dropdown menu -->
    <div
      v-if="isOpen"
      class="absolute z-50 w-full mt-1 bg-white border border-light-border rounded-lg shadow-lg max-h-60 overflow-auto"
    >
      <ul class="py-1">
        <li
          v-for="(option, index) in options"
          :key="index"
          @click="selectOption(option)"
          class="px-4 py-2 text-sm cursor-pointer transition-colors"
          :class="[
            option.disabled ? 'text-gray-400 cursor-not-allowed bg-gray-50' : 'text-slate-700 hover:bg-blue-50 hover:text-blue-600',
            modelValue === option.value ? 'bg-blue-50 text-blue-600 font-medium' : ''
          ]"
        >
          {{ option.label }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  options: {
    type: Array as () => Array<{ label: string, value: any, disabled?: boolean }>,
    required: true
  },
  placeholder: {
    type: String,
    default: 'Vui lòng chọn'
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const isOpen = ref(false)
const container = ref<HTMLElement | null>(null)

const selectedOption = computed(() => {
  return props.options.find(opt => opt.value === props.modelValue)
})

const selectOption = (option: any) => {
  if (option.disabled) return
  emit('update:modelValue', option.value)
  emit('change', option.value)
  isOpen.value = false
}

// Click outside to close
const handleClickOutside = (event: MouseEvent) => {
  if (container.value && !container.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

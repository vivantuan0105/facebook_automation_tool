<template>
  <div class="relative" ref="container">
    <div
      @click="isOpen = !isOpen"
      class="w-full text-sm bg-light-bg border border-light-border rounded-lg p-2.5 cursor-pointer flex justify-between items-center transition-colors"
      :class="[isOpen ? 'ring-2 ring-primary border-primary' : 'hover:border-primary']"
    >
      <span :class="{'text-light-muted': !selectedLabels.length}" class="truncate">
        {{ selectedLabels.length ? selectedLabels.join(', ') : placeholder }}
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
          @click.stop="toggleAll"
          class="px-4 py-2 text-sm cursor-pointer transition-colors text-slate-700 hover:bg-slate-50 flex items-center font-medium border-b border-gray-100"
        >
          <input type="checkbox" :checked="isAllSelected" class="mr-3 rounded text-primary focus:ring-primary h-4 w-4" />
          Chọn tất cả
        </li>
        <li
          v-for="(option, index) in options"
          :key="index"
          @click.stop="toggleOption(option)"
          class="px-4 py-2 text-sm cursor-pointer transition-colors flex items-center"
          :class="[
            option.disabled ? 'text-gray-400 cursor-not-allowed bg-gray-50' : 'text-slate-700 hover:bg-blue-50',
            isSelected(option) ? 'bg-blue-50 text-blue-600 font-medium' : ''
          ]"
        >
          <input type="checkbox" :checked="isSelected(option)" :disabled="option.disabled" class="mr-3 rounded text-primary focus:ring-primary h-4 w-4" />
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
    type: Array as () => any[],
    default: () => []
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

const selectedLabels = computed(() => {
  return props.options.filter(opt => props.modelValue.includes(opt.value)).map(opt => opt.label)
})

const isAllSelected = computed(() => {
  const activeOptions = props.options.filter(opt => !opt.disabled)
  if (activeOptions.length === 0) return false
  return activeOptions.every(opt => props.modelValue.includes(opt.value))
})

const toggleAll = () => {
  const activeOptions = props.options.filter(opt => !opt.disabled)
  if (isAllSelected.value) {
    emit('update:modelValue', [])
    emit('change', [])
  } else {
    const allValues = activeOptions.map(opt => opt.value)
    emit('update:modelValue', allValues)
    emit('change', allValues)
  }
}

const isSelected = (option: any) => {
  return props.modelValue.includes(option.value)
}

const toggleOption = (option: any) => {
  if (option.disabled) return
  let newValue = [...props.modelValue]
  if (isSelected(option)) {
    newValue = newValue.filter(v => v !== option.value)
  } else {
    newValue.push(option.value)
  }
  emit('update:modelValue', newValue)
  emit('change', newValue)
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

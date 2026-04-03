import { defineStore } from 'pinia'
import { ref } from 'vue'
import { tournamentsApi, participantsApi, encountersApi, usersApi } from '@/api'

export const useTournamentStore = defineStore('tournament', () => {
  const owning = ref([])
  const participating = ref([])
  const current = ref(null)
  const participants = ref([])
  const encounters = ref([])
  const myEncounters = ref([])
  const scoreboard = ref([])
  const users = ref([])
  const loading = ref(false)

  async function fetchAll() {
    loading.value = true
    try {
      const { data } = await tournamentsApi.list()
      owning.value = data.owning || []
      participating.value = data.participating || []
    } finally {
      loading.value = false
    }
  }

  async function fetchOne(id) {
    loading.value = true
    try {
      const { data } = await tournamentsApi.get(id)
      current.value = data
      return data
    } finally {
      loading.value = false
    }
  }

  async function fetchParticipants(tid) {
    const { data } = await participantsApi.list(tid)
    data.sort((a, b) => a.visible_name.toLocaleLowerCase().localeCompare(b.visible_name.toLocaleLowerCase()))
    participants.value = data
    return data
  }

  async function fetchEncounters(tid, userId) {
    const { data } = await encountersApi.list(tid, userId)
    if (userId) myEncounters.value = data
    else encounters.value = data
    return data
  }

  async function fetchScoreboard(tid) {
    const { data } = await tournamentsApi.scoreboard(tid)
    scoreboard.value = data
    return data
  }

  async function fetchUsers() {
    try {
      const { data } = await usersApi.list()
      users.value = data
    } catch {}
  }

  async function createUser(payload) {
    const {data}  = await usersApi.create(payload)
    return data
  }

  async function createTournament(payload) {
    const { data } = await tournamentsApi.create(payload)
    return data
  }

  async function updateTournament(id, payload) {
    const { data } = await tournamentsApi.update(id, payload)
    if (current.value?.id === id) current.value = data
    return data
  }

  async function deleteTournament(id) {
    await tournamentsApi.delete(id)
    owning.value = owning.value.filter(t => t.id !== id)
  }

  async function resetTournament(id) {
    await tournamentsApi.reset(id)
  }

  async function addParticipants(tid, ids) {
    await participantsApi.add(tid, ids)
    await fetchParticipants(tid)
  }

  async function removeParticipant(tid, id) {
    await participantsApi.remove(tid, id)
    await fetchParticipants(tid)
  }

  async function createEncounters(tid, list) {
    await encountersApi.create(tid, list)
    await fetchEncounters(tid)
  }

  async function updateScore(tid, encounter) {
    await encountersApi.updateScore(tid, encounter)
    await fetchEncounters(tid)
  }

  async function deleteScore(tid, data) {
    await encountersApi.deleteScore(tid, data)
    await fetchEncounters(tid)
  }

  return {
    owning, participating, current, participants, encounters, myEncounters, scoreboard, users, loading,
    fetchAll, fetchOne, fetchParticipants, fetchEncounters, fetchScoreboard, fetchUsers,
    createTournament, updateTournament, deleteTournament, resetTournament,
    addParticipants, removeParticipant, createEncounters, updateScore, deleteScore,
    createUser,
  }
})

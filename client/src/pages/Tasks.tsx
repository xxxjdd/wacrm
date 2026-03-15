import React, { useEffect, useState } from 'react'
import { tasksApi } from '../api/client'
import { Plus, Play, Edit2, Trash2, Calendar, Clock } from 'lucide-react'

export default function Tasks() {
  const [tasks, setTasks] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadTasks()
  }, [])

  const loadTasks = async () => {
    try {
      const { data } = await tasksApi.list()
      setTasks(data.tasks || [])
    } catch (error) {
      console.error('Failed to load tasks:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleRun = async (id: number) => {
    try {
      await tasksApi.run(id)
      alert('Task started')
    } catch (error) {
      alert('Failed to run task')
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure?')) return
    try {
      await tasksApi.delete(id)
      loadTasks()
    } catch (error) {
      alert('Failed to delete task')
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed': return 'bg-green-100 text-green-700'
      case 'running': return 'bg-blue-100 text-blue-700'
      case 'pending': return 'bg-yellow-100 text-yellow-700'
      case 'failed': return 'bg-red-100 text-red-700'
      default: return 'bg-gray-100 text-gray-700'
    }
  }

  return (
    <div className="p-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Scheduled Tasks</h1>
        <button
          onClick={() => alert('Create task modal - TODO')}
          className="flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
        >
          <Plus size={20} />
          New Task
        </button>
      </div>

      {loading ? (
        <div className="flex justify-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      ) : tasks.length === 0 ? (
        <div className="bg-white rounded-lg shadow p-8 text-center text-gray-500">
          No scheduled tasks yet.
        </div>
      ) : (
        <div className="space-y-4">
          {tasks.map((task) => (
            <div key={task.id} className="bg-white rounded-lg shadow p-4">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-semibold">{task.name}</h3>
                  <p className="text-sm text-gray-500">
                    Scheduled: {task.scheduled_at && new Date(task.scheduled_at).toLocaleString()}
                  </p>
                </div>
                <span className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(task.status)}`}>
                  {task.status}
                </span>
              </div>
              <div className="flex items-center gap-4 mt-3 text-sm text-gray-500">
                <span className="flex items-center gap-1">
                  <Calendar size={14} />
                  {task.repeat_type || 'Once'}
                </span>
                <span className="flex items-center gap-1">
                  <Clock size={14} />
                  Run count: {task.run_count}
                </span>
              </div>
              <div className="flex gap-2 mt-3">
                <button
                  onClick={() => handleRun(task.id)}
                  disabled={task.status === 'running'}
                  className="flex items-center gap-1 px-3 py-1 bg-blue-100 text-blue-700 rounded hover:bg-blue-200 disabled:opacity-50"
                >
                  <Play size={14} />
                  Run Now
                </button>
                <button
                  onClick={() => alert('Edit - TODO')}
                  className="p-1 text-gray-600 hover:text-blue-600"
                >
                  <Edit2 size={16} />
                </button>
                <button
                  onClick={() => handleDelete(task.id)}
                  className="p-1 text-gray-600 hover:text-red-600"
                >
                  <Trash2 size={16} />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

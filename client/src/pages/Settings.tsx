import React from 'react'
import { useAuthStore } from '../store/auth'
import { LogOut, User, Shield, Bell } from 'lucide-react'

export default function Settings() {
  const { user, logout } = useAuthStore()

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold text-gray-800 mb-6">Settings</h1>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Profile Section */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center gap-3 mb-4">
            <User className="text-blue-600" size={24} />
            <h2 className="text-lg font-semibold">Profile</h2>
          </div>
          
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Username
              </label>
              <input
                type="text"
                value={user?.username || ''}
                readOnly
                className="w-full px-3 py-2 border rounded-lg bg-gray-50"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Email
              </label>
              <input
                type="email"
                value={user?.email || ''}
                readOnly
                className="w-full px-3 py-2 border rounded-lg bg-gray-50"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Role
              </label>
              <input
                type="text"
                value={user?.role || ''}
                readOnly
                className="w-full px-3 py-2 border rounded-lg bg-gray-50"
              />
            </div>
          </div>
        </div>

        {/* Security Section */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center gap-3 mb-4">
            <Shield className="text-green-600" size={24} />
            <h2 className="text-lg font-semibold">Security</h2>
          </div>
          
          <div className="space-y-4">
            <button className="w-full text-left px-4 py-3 rounded-lg border hover:bg-gray-50 transition-colors">
              <p className="font-medium">Change Password</p>
              <p className="text-sm text-gray-500">Update your account password</p>
            </button>
            
            <button className="w-full text-left px-4 py-3 rounded-lg border hover:bg-gray-50 transition-colors">
              <p className="font-medium">Two-Factor Authentication</p>
              <p className="text-sm text-gray-500">Add extra security to your account</p>
            </button>
          </div>
        </div>

        {/* Notifications Section */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center gap-3 mb-4">
            <Bell className="text-orange-600" size={24} />
            <h2 className="text-lg font-semibold">Notifications</h2>
          </div>
          
          <div className="space-y-3">
            <label className="flex items-center gap-3">
              <input type="checkbox" className="w-4 h-4 text-blue-600 rounded" defaultChecked />
              <span>New message notifications</span>
            </label>
            <label className="flex items-center gap-3">
              <input type="checkbox" className="w-4 h-4 text-blue-600 rounded" defaultChecked />
              <span>Task completion notifications</span>
            </label>
            <label className="flex items-center gap-3">
              <input type="checkbox" className="w-4 h-4 text-blue-600 rounded" />
              <span>Email notifications</span>
            </label>
          </div>
        </div>

        {/* Logout Section */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center gap-3 mb-4">
            <LogOut className="text-red-600" size={24} />
            <h2 className="text-lg font-semibold">Account</h2>
          </div>
          
          <button
            onClick={logout}
            className="w-full bg-red-100 text-red-700 py-3 rounded-lg hover:bg-red-200 transition-colors font-medium"
          >
            Logout
          </button>
        </div>
      </div>
    </div>
  )
}

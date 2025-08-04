import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation } from 'react-query';
import { toast } from 'react-hot-toast';
import { 
  Package, 
  Play, 
  Trash2, 
  Edit,
  ExternalLink,
  Activity,
  Clock,
  GitBranch
} from 'lucide-react';
import { applicationApi } from '../services/api';
import { Application } from '../types';

const ApplicationDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data: application, isLoading, error } = useQuery(
    ['application', id],
    () => applicationApi.getApplication(id!),
    { enabled: !!id }
  );

  const deployMutation = useMutation(
    () => applicationApi.deployApplication(id!),
    {
      onSuccess: () => {
        toast.success('Deployment started!');
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.error || 'Failed to start deployment');
      },
    }
  );

  const deleteMutation = useMutation(
    () => applicationApi.deleteApplication(id!),
    {
      onSuccess: () => {
        toast.success('Application deleted!');
        navigate('/applications');
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.error || 'Failed to delete application');
      },
    }
  );

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="spinner w-8 h-8"></div>
      </div>
    );
  }

  if (error || !application) {
    return (
      <div className="text-center py-12">
        <Package className="mx-auto h-12 w-12 text-red-400" />
        <h3 className="mt-2 text-sm font-medium text-gray-900">Application not found</h3>
        <p className="mt-1 text-sm text-gray-500">The application you're looking for doesn't exist.</p>
      </div>
    );
  }

  const handleDelete = () => {
    if (window.confirm('Are you sure you want to delete this application? This action cannot be undone.')) {
      deleteMutation.mutate();
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-start">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">{application.name}</h1>
          <p className="mt-1 text-sm text-gray-500">
            {application.description || 'No description provided'}
          </p>
        </div>
        <div className="flex space-x-3">
          <button
            onClick={() => deployMutation.mutate()}
            disabled={deployMutation.isLoading}
            className="btn-success"
          >
            <Play className="h-4 w-4 mr-2" />
            {deployMutation.isLoading ? 'Deploying...' : 'Deploy'}
          </button>
          <button
            onClick={() => navigate(`/applications/${id}/edit`)}
            className="btn-secondary"
          >
            <Edit className="h-4 w-4 mr-2" />
            Edit
          </button>
          <button
            onClick={handleDelete}
            disabled={deleteMutation.isLoading}
            className="btn-danger"
          >
            <Trash2 className="h-4 w-4 mr-2" />
            {deleteMutation.isLoading ? 'Deleting...' : 'Delete'}
          </button>
        </div>
      </div>

      {/* Status and Info */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Main Info */}
        <div className="lg:col-span-2 space-y-6">
          {/* Status Card */}
          <div className="card">
            <div className="card-header">
              <h3 className="text-lg font-medium text-gray-900">Status</h3>
            </div>
            <div className="card-body">
              <div className="flex items-center justify-between">
                <div className="flex items-center">
                  <span className={`status-indicator ${getStatusClass(application.status)}`}></span>
                  <span className="ml-2 text-sm font-medium text-gray-900 capitalize">
                    {application.status}
                  </span>
                </div>
                <span className={`badge ${getStatusColor(application.status)}`}>
                  {application.status}
                </span>
              </div>
            </div>
          </div>

          {/* Application Details */}
          <div className="card">
            <div className="card-header">
              <h3 className="text-lg font-medium text-gray-900">Application Details</h3>
            </div>
            <div className="card-body">
              <dl className="grid grid-cols-1 gap-x-4 gap-y-6 sm:grid-cols-2">
                <div>
                  <dt className="text-sm font-medium text-gray-500">Environment</dt>
                  <dd className="mt-1 text-sm text-gray-900">
                    <span className={`badge ${getEnvironmentColor(application.environment)}`}>
                      {application.environment}
                    </span>
                  </dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500">Port</dt>
                  <dd className="mt-1 text-sm text-gray-900">{application.port}</dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500">Replicas</dt>
                  <dd className="mt-1 text-sm text-gray-900">{application.replicas}</dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500">Created</dt>
                  <dd className="mt-1 text-sm text-gray-900">
                    {new Date(application.created_at).toLocaleDateString()}
                  </dd>
                </div>
              </dl>
            </div>
          </div>

          {/* Git Repository */}
          <div className="card">
            <div className="card-header">
              <h3 className="text-lg font-medium text-gray-900">Git Repository</h3>
            </div>
            <div className="card-body">
              <dl className="grid grid-cols-1 gap-x-4 gap-y-6 sm:grid-cols-2">
                <div>
                  <dt className="text-sm font-medium text-gray-500">Repository URL</dt>
                  <dd className="mt-1 text-sm text-gray-900">
                    <a 
                      href={application.git_url} 
                      target="_blank" 
                      rel="noopener noreferrer"
                      className="text-primary-600 hover:text-primary-500 flex items-center"
                    >
                      {application.git_url}
                      <ExternalLink className="h-3 w-3 ml-1" />
                    </a>
                  </dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500">Branch</dt>
                  <dd className="mt-1 text-sm text-gray-900 flex items-center">
                    <GitBranch className="h-4 w-4 mr-1" />
                    {application.branch}
                  </dd>
                </div>
              </dl>
            </div>
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Quick Actions */}
          <div className="card">
            <div className="card-header">
              <h3 className="text-lg font-medium text-gray-900">Quick Actions</h3>
            </div>
            <div className="card-body space-y-3">
              <button
                onClick={() => navigate(`/applications/${id}/logs`)}
                className="w-full btn-secondary justify-start"
              >
                <Activity className="h-4 w-4 mr-2" />
                View Logs
              </button>
              <button
                onClick={() => navigate(`/applications/${id}/metrics`)}
                className="w-full btn-secondary justify-start"
              >
                <Activity className="h-4 w-4 mr-2" />
                View Metrics
              </button>
              {application.deployment_url && (
                <a
                  href={application.deployment_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="w-full btn-secondary justify-start"
                >
                  <ExternalLink className="h-4 w-4 mr-2" />
                  Open Application
                </a>
              )}
            </div>
          </div>

          {/* Image Info */}
          {application.image_url && (
            <div className="card">
              <div className="card-header">
                <h3 className="text-lg font-medium text-gray-900">Docker Image</h3>
              </div>
              <div className="card-body">
                <p className="text-sm text-gray-900 break-all">
                  {application.image_url}
                </p>
              </div>
            </div>
          )}

          {/* Last Updated */}
          <div className="card">
            <div className="card-header">
              <h3 className="text-lg font-medium text-gray-900">Last Updated</h3>
            </div>
            <div className="card-body">
              <div className="flex items-center text-sm text-gray-900">
                <Clock className="h-4 w-4 mr-2" />
                {new Date(application.updated_at).toLocaleString()}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

// Helper functions
const getStatusColor = (status: string) => {
  switch (status) {
    case 'running':
      return 'bg-green-100 text-green-800';
    case 'building':
    case 'deploying':
      return 'bg-yellow-100 text-yellow-800';
    case 'failed':
      return 'bg-red-100 text-red-800';
    case 'stopped':
      return 'bg-gray-100 text-gray-800';
    default:
      return 'bg-blue-100 text-blue-800';
  }
};

const getStatusClass = (status: string) => {
  switch (status) {
    case 'running':
      return 'status-running';
    case 'building':
    case 'deploying':
      return 'status-building';
    case 'failed':
      return 'status-failed';
    case 'stopped':
      return 'status-stopped';
    default:
      return 'status-pending';
  }
};

const getEnvironmentColor = (environment: string) => {
  switch (environment) {
    case 'production':
      return 'bg-red-100 text-red-800';
    case 'staging':
      return 'bg-orange-100 text-orange-800';
    default:
      return 'bg-blue-100 text-blue-800';
  }
};

export default ApplicationDetail; 
import React from 'react';
import { useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { useMutation } from 'react-query';
import { toast } from 'react-hot-toast';
import { applicationApi } from '../services/api';
import { ApplicationCreate } from '../types';

const CreateApplication: React.FC = () => {
  const navigate = useNavigate();
  const { register, handleSubmit, formState: { errors } } = useForm<ApplicationCreate>();

  const createMutation = useMutation(applicationApi.createApplication, {
    onSuccess: (data) => {
      toast.success('Application created successfully!');
      navigate(`/applications/${data.id}`);
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to create application');
    },
  });

  const onSubmit = (data: ApplicationCreate) => {
    createMutation.mutate(data);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Create New Application</h1>
        <p className="mt-1 text-sm text-gray-500">
          Deploy your application to the platform
        </p>
      </div>

      {/* Form */}
      <div className="card">
        <div className="card-body">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
              <div>
                <label htmlFor="name" className="block text-sm font-medium text-gray-700">
                  Application Name *
                </label>
                <input
                  id="name"
                  type="text"
                  {...register('name', { 
                    required: 'Application name is required',
                    minLength: {
                      value: 1,
                      message: 'Name must be at least 1 character'
                    },
                    maxLength: {
                      value: 100,
                      message: 'Name must be less than 100 characters'
                    }
                  })}
                  className="input mt-1"
                  placeholder="my-awesome-app"
                />
                {errors.name && (
                  <p className="mt-1 text-sm text-red-600">{errors.name.message}</p>
                )}
              </div>

              <div>
                <label htmlFor="environment" className="block text-sm font-medium text-gray-700">
                  Environment *
                </label>
                <select
                  id="environment"
                  {...register('environment', { required: 'Environment is required' })}
                  className="input mt-1"
                >
                  <option value="development">Development</option>
                  <option value="staging">Staging</option>
                  <option value="production">Production</option>
                </select>
                {errors.environment && (
                  <p className="mt-1 text-sm text-red-600">{errors.environment.message}</p>
                )}
              </div>
            </div>

            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700">
                Description
              </label>
              <textarea
                id="description"
                rows={3}
                {...register('description')}
                className="input mt-1"
                placeholder="Describe your application..."
              />
            </div>

            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
              <div>
                <label htmlFor="git_url" className="block text-sm font-medium text-gray-700">
                  Git Repository URL *
                </label>
                <input
                  id="git_url"
                  type="url"
                  {...register('git_url', { 
                    required: 'Git URL is required',
                    pattern: {
                      value: /^https?:\/\/.+/,
                      message: 'Please enter a valid URL'
                    }
                  })}
                  className="input mt-1"
                  placeholder="https://github.com/username/repo.git"
                />
                {errors.git_url && (
                  <p className="mt-1 text-sm text-red-600">{errors.git_url.message}</p>
                )}
              </div>

              <div>
                <label htmlFor="branch" className="block text-sm font-medium text-gray-700">
                  Branch *
                </label>
                <input
                  id="branch"
                  type="text"
                  {...register('branch', { 
                    required: 'Branch is required'
                  })}
                  className="input mt-1"
                  placeholder="main"
                />
                {errors.branch && (
                  <p className="mt-1 text-sm text-red-600">{errors.branch.message}</p>
                )}
              </div>
            </div>

            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
              <div>
                <label htmlFor="port" className="block text-sm font-medium text-gray-700">
                  Port *
                </label>
                <input
                  id="port"
                  type="number"
                  {...register('port', { 
                    required: 'Port is required',
                    min: {
                      value: 1,
                      message: 'Port must be between 1 and 65535'
                    },
                    max: {
                      value: 65535,
                      message: 'Port must be between 1 and 65535'
                    }
                  })}
                  className="input mt-1"
                  placeholder="8080"
                />
                {errors.port && (
                  <p className="mt-1 text-sm text-red-600">{errors.port.message}</p>
                )}
              </div>

              <div>
                <label htmlFor="replicas" className="block text-sm font-medium text-gray-700">
                  Replicas *
                </label>
                <input
                  id="replicas"
                  type="number"
                  {...register('replicas', { 
                    required: 'Replicas is required',
                    min: {
                      value: 1,
                      message: 'Replicas must be between 1 and 10'
                    },
                    max: {
                      value: 10,
                      message: 'Replicas must be between 1 and 10'
                    }
                  })}
                  className="input mt-1"
                  placeholder="1"
                />
                {errors.replicas && (
                  <p className="mt-1 text-sm text-red-600">{errors.replicas.message}</p>
                )}
              </div>
            </div>

            <div className="flex justify-end space-x-3">
              <button
                type="button"
                onClick={() => navigate('/applications')}
                className="btn-secondary"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={createMutation.isLoading}
                className="btn-primary"
              >
                {createMutation.isLoading ? (
                  <div className="flex items-center">
                    <div className="spinner w-4 h-4 mr-2"></div>
                    Creating...
                  </div>
                ) : (
                  'Create Application'
                )}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default CreateApplication; 
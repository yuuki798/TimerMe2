// TaskList.tsx
import React from 'react';

interface Task {
  id: number;
  name: string;
  duration: number;
  is_completed: boolean;
  start_time: string;
  status: string;
  total_time: number;
}

interface TaskListProps {
  tasks: Task[];
  handleStartTask: (id: number) => void;
  handlePauseTask: (id: number) => void;
  handleCompleteTask: (id: number) => void;
  handleDeleteTask: (id: number) => void;
  handleResetTask: (id: number) => void; // 新增重置任务处理函数
}

const formatTimeWithUnit = (seconds: number) => {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;

  return `${hours}h ${minutes}m ${secs}s`;
};

const TaskList: React.FC<TaskListProps> = ({
  tasks,
  handleStartTask,
  handlePauseTask,
  handleCompleteTask,
  handleDeleteTask,
  handleResetTask, // 新增重置任务处理函数
}) => {
  return (
    <ul className={'list-none p-0'}>
      {tasks.map((task: Task) => {
        const progress = (task.duration / task.total_time) * 100;
        return (
          <li
            key={task.id}
            className={'border-b p-4 flex justify-between items-center'}
          >
            <span className={'flex-grow w-2/3'}>
              {task.name} —— {formatTimeWithUnit(task.duration)} /{' '}
              {formatTimeWithUnit(task.total_time)} ({progress.toFixed(2)}%)
            </span>

            <span>{task.status}</span>

            <div className='progress-bar w-full bg-gray-200 rounded-full h-4 overflow-hidden'>
              <div
                className='progress bg-blue-500 h-full'
                style={{ width: `${progress}%` }}
              ></div>
            </div>

            <div className='flex flex-col items-end space-y-2'>
              <button
                onClick={() => handleStartTask(task.id)}
                className={'p-2 rounded bg-green-500 text-white w-20'}
              >
                Start
              </button>
              <button
                onClick={() => handlePauseTask(task.id)}
                className={'p-2 rounded bg-yellow-500 text-white w-20'}
              >
                Pause
              </button>
              <button
                onClick={() => handleDeleteTask(task.id)}
                className={'p-2 rounded bg-red-500 text-white w-20'}
              >
                Delete
              </button>
              <button
                onClick={() => handleResetTask(task.id)} // 新增重置任务按钮
                className={'p-2 rounded bg-gray-500 text-white w-20'}
              >
                Reset
              </button>
            </div>
          </li>
        );
      })}
    </ul>
  );
};

export default TaskList;

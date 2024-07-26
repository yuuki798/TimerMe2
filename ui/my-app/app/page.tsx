'use client';
import React, { useState, useEffect } from 'react';
import {
  getTasks,
  createTask,
  updateTask,
  deleteTask,
  startTask,
  pauseTask,
  completeTask,
} from '@/utils/api';
import TaskList from '@/components/TaskList';
import { Simulate } from 'react-dom/test-utils';
import pause = Simulate.pause;

interface Task {
  id: number;
  name: string;
  duration: number;
  is_completed: boolean;
  start_time: string;
  status: string;
  total_time: number;
}

export default function Home() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [newTaskName, setNewTaskName] = useState('');
  const [newTotalTime, setNewTotalTime] = useState<number>(0);
  const [taskTimers, setTaskTimers] = useState<{
    [key: number]: NodeJS.Timeout;
  }>({});
  const [timeUnit, setTimeUnit] = useState('seconds');

  useEffect(() => {
    fetchTasks();
  }, []);

  const startTaskTimer = (id: number) => {
    const timer = setInterval(() => {
      setTasks((prevTasks) => {
        return prevTasks.map((task) => {
          if (task.duration + 1 >= task.total_time) {
            return { ...task, duration: task.total_time, status: 'completed' };
          }

          if (task.id === id && task.status === 'started') {
            return { ...task, duration: task.duration + 1 };
          }
          return task;
        });
      });
    }, 1000);
    setTaskTimers((prevTimers) => ({ ...prevTimers, [id]: timer }));
  };

  const stopTaskTimer = (id: number) => {
    if (taskTimers[id]) {
      clearInterval(taskTimers[id]);
      setTaskTimers((timers) => {
        const newTimers = { ...timers };
        delete newTimers[id];
        return newTimers;
      });
    }
  };
  const fetchTasks = async () => {
    const tasks = await getTasks();
    setTasks(tasks);
  };

  const handleCreateTask = async () => {
    let totalTimeInSeconds = newTotalTime;

    // 根据选择的计时单位转换为秒
    switch (timeUnit) {
      case 'minutes':
        totalTimeInSeconds = newTotalTime * 60;
        break;
      case 'hours':
        totalTimeInSeconds = newTotalTime * 3600;
        break;
      default:
        break;
    }

    const newTask = await createTask(newTaskName, totalTimeInSeconds);
    setTasks([...tasks, newTask]);
    setNewTaskName('');
    setNewTotalTime(0);
    setTimeUnit('seconds');
  };

  const handleUpdateTask = async (id: number, updatedTask: any) => {
    const task = await updateTask(id, updatedTask);
    setTasks(tasks.map((t: any) => (t.id === id ? task : t)));
  };
  const handleDeleteTask = async (id: number) => {
    await deleteTask(id);
    setTasks(tasks.filter((t: any) => t.id !== id));
  };
  const handleStartTask = async (id: number) => {
    const task = await startTask(id);
    setTasks(tasks.map((t: any) => (t.id == id ? task : t)));
    startTaskTimer(id);
  };

  const handlePauseTask = async (id: number) => {
    const task = await pauseTask(id);
    setTasks(tasks.map((t: any) => (t.id == id ? task : t)));
    stopTaskTimer(id);
  };

  const handleCompleteTask = async (id: number) => {
    const task = await completeTask(id);
    setTasks(tasks.map((t: any) => (t.id == id ? task : t)));
    stopTaskTimer(id);
  };
  return (
    <div className={'container mx-auto p-4'}>
      <h1 className={'text-2xl font-bold mb-4'}>Welcome to TimerMe!</h1>
      <div className={'mb-4'}>
        <input
          type={'text'}
          value={newTaskName}
          onChange={(e) => setNewTaskName(e.target.value)}
          placeholder={'New Task Name'}
          className={'border p-2 mr-2'}
        />
        <input
          type={'number'}
          value={newTotalTime}
          onChange={(e) => setNewTotalTime(parseInt(e.target.value))}
          placeholder={'Duration(seconds)'}
          className={'border p-2 mr-2'}
        />

        <select
          value={timeUnit}
          onChange={(e) => setTimeUnit(e.target.value)}
          className='border p-2 mr-2'
        >
          <option value='seconds'>Seconds</option>
          <option value='minutes'>Minutes</option>
          <option value='hours'>Hours</option>
        </select>

        <button
          onClick={handleCreateTask}
          className={'bg-blue-500 text-white p-2 rounded'}
        >
          Add Task
        </button>
        <TaskList
          tasks={tasks}
          handleStartTask={handleStartTask}
          handlePauseTask={handlePauseTask}
          handleCompleteTask={handleCompleteTask}
          handleDeleteTask={handleDeleteTask}
        />
      </div>
    </div>
  );
}

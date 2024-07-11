import { useEffect, useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

function App() {
  const [todoList, setTodoList] = useState([])
  const [title, setTitle] = useState()
  const [editValue, setEditValue] = useState()

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    try {
      const response = await fetch('http://localhost:8091/api/v1/todos');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const result = await response.json();
      const newResult = result.map((e) => ({ ...e, isEdit: false }));
      setTodoList(newResult)

    } catch (err) {
      console.log(err);
    }
  }

  const handleUpdateTodo = async (id) => {
    const payload = { Title: editValue }
    try {
      const response = await fetch(`http://localhost:8091/api/v1/todos/${id}/actions/title`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      fetchData()
    } catch (err) {
      console.log(err);
    }
  }

  const handleOpenEdit = (id) => {
    const updatedTodos = todoList.map((todo) =>
      todo.ID === id ? { ...todo, isEdit: true } : { ...todo, isEdit: false }
    );
    setTodoList(updatedTodos);
  }

  const handleCloseEditt = (id) => {
    const updatedTodos = todoList.map((todo) => ({ ...todo, isEdit: false }));
    setTodoList(updatedTodos);
  }

  const handleDeleteTodo = async (id) => {
    try {
      const response = await fetch(`http://localhost:8091/api/v1/todos/${id}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      fetchData()
    } catch (err) {
      console.log(err);
    }
  }

  const handleAddTodo = async () => {
    const payload = { Title: title, Status: 'inactive' }
    try {
      const response = await fetch(`http://localhost:8091/api/v1/todos`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      fetchData()
      setTitle()
    } catch (err) {
      console.log(err);
    }
  }

  const handleChangeStatus = async (id, status) => {
    const payload = { Status: status }
    try {
      const response = await fetch(`http://localhost:8091/api/v1/todos/${id}/actions/status`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      fetchData()
    } catch (err) {
      console.log(err);
    }
  }

  return (
    <div className='bg-white rounded-2xl p-8 w-[500px]'>
      <div className='flex gap-4 items-center justify-between mb-5'>
        <button className="btn btn-active btn-primary" onClick={() => { handleAddTodo() }}>Add</button>
        <label className="input input-bordered flex items-center gap-6 w-full">
          What to do?
          <input type="text" className="grow" placeholder="write..." value={title} onChange={(e) => setTitle(e.target.value)} />
        </label>
      </div>


      <ul className="list-decimal list-inside">
        {todoList.map((e) => {
          return (
            e.isEdit ?
              <div className='flex items-center gap-3 justify-between my-4' key={e.ID}>
                <input type="text" className="input input-bordered w-full max-w-xs" value={editValue} onChange={(e) => setEditValue(e.target.value)} />
                <div className='flex items-center gap-3'>
                  <button className="btn btn-active btn-ghost" onClick={() => handleUpdateTodo(e.ID)}>Save</button>
                  <button className="btn btn-square btn-error text-white"
                    onClick={() => {
                      setEditValue()
                      handleCloseEditt(e.ID)
                    }}
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-6 w-6"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor">
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
              </div>

              : <div className='flex align-center justify-between py-3' key={e.ID}>
                <li className='text-left ps-3' onClick={() => {
                  setEditValue(e.Title)
                  handleOpenEdit(e.ID)
                }} >{e.Title}</li>
                <div className='flex items-center gap-3'>
                  {e.Status == 'active' ?
                    <div className="badge badge-success text-white" onClick={() => handleChangeStatus(e.ID, 'inactive')}>
                      Done
                    </div> :
                    <div className="badge badge-outline" onClick={() => handleChangeStatus(e.ID, 'active')}>
                      To Do
                    </div>
                  }
                  <button className="btn btn-xs btn-error text-white btn-circle" onClick={() => handleDeleteTodo(e.ID)}>
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-6 w-6"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor">
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>

              </div>
          )
        })}
      </ul>
    </div>
  )
}


export default App

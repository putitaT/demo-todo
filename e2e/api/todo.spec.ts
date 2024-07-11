import { test, expect } from '@playwright/test';

test.describe('Get all todos list', () => {
    test('should respond all todo items that adding to /api/v1/todos', async ({ request }) => {
        const postRes1 = await request.post('http://localhost:8091/api/v1/todos',
            {
                data: {
                    Title: 'Tailwind CSS',
                    Status: 'active'
                }
            }
        )
        const postRes2 = await request.post('http://localhost:8091/api/v1/todos',
            {
                data: {
                    Title: 'Bootstap',
                    Status: 'inactive'
                }
            }
        )
        const resp = await request.get('http://localhost:8091/api/v1/todos')

        expect(resp.ok()).toBeTruthy()
        expect(await resp.json()).toEqual(
            expect.arrayContaining([
                {
                    ID: expect.any(Number),
                    Title: 'Tailwind CSS',
                    Status: 'active'
                },
                {
                    ID: expect.any(Number),
                    Title: 'Bootstap',
                    Status: 'inactive'
                }
            ]
            )
        )

        const idTodo1 = await postRes1.json()
        const idTodo2 = await postRes2.json()

        await request.delete('http://localhost:8091/api/v1/todos/' + String(idTodo1))
        await request.delete('http://localhost:8091/api/v1/todos/' + String(idTodo2))
    });
})

test.describe('Get todo by ID', () => {
    test('should respond with one todo when add todo /api/v1/todos/:id', async ({
        request,
    }) => {
        const postResponse = await request.post('http://localhost:8091/api/v1/todos',
            {
                data: {
                    Title: 'Angular',
                    Status: 'active'
                }
            }
        )
        const newTodo = await postResponse.json()

        const response = await request.get('http://localhost:8091/api/v1/todos/' + String(newTodo))

        expect(response.ok()).toBeTruthy()
        expect(await response.json()).toEqual(
            expect.objectContaining({
                ID: expect.any(Number),
                Title: 'Angular',
                Status: 'active'
            })
        )

        await request.delete('http://localhost:8091/api/v1/todos/' + String(newTodo))
    })
})





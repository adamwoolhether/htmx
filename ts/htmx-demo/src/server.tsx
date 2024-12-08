import {type Context, Hono} from 'hono';
import {serveStatic} from 'hono/bun';

const app = new Hono();

app.use('/*', serveStatic({root: './public'}));
// app.get('/version', (c: Context) => {
//   return c.text(Bun.version)
// })
app.get('/table-rows', (c: Context) => {
  const sortedDogs = Array.from(dogs.values()).sort((a, b) =>
    a.name.localeCompare(b.name)
  );

  return c.html(<>{sortedDogs.map(dog => dogRow(dog))}</>);
});
app.post('/dog', async (c: Context) => {
  const formData = await c.req.formData();
  const name = (formData.get('name') as string) || '';
  const breed = (formData.get('breed') as string) || '';
  const dog = addDog(name, breed);
  return c.html(dogRow(dog), 201);
});
app.put('/select/:id', (c: Context) => {
  selectId = c.req.param('id');
  c.header('HX-Trigger', 'selection-change')
  return c.body(null);
})
app.put('/dog/:id', async (c: Context) => {
  const id = c.req.param('id');
  const formData = await c.req.formData();
  const name = (formData.get('name') as string) || '';
  const breed = (formData.get('breed') as string) || '';
  const updatedDog = {id, name, breed};
  dogs.set(id, updatedDog);
  selectId = '';
  c.header('HX-Trigger', 'selection-change')

  return c.html(dogRow(updatedDog, true))
})
app.put('/deselect', (c: Context) => {
  selectId = '';
  c.header('HX-Trigger', 'selection-change')
  return c.body(null);
})
app.delete('/dog/:id', async (c: Context) => {
  const id = c.req.param('id');
  dogs.delete(id);
  return c.body(null);
});

app.get('/form', (c: Context) => {
  const attrs: {[key: string]: string} = {
    'hx-on:htmx:after-request': 'this.reset()'
  };

  if (selectId) {
    // Update an existing row.
    attrs['hx-put'] = '/dog/' + selectId;
  } else {
    // Add new row.
    attrs['hx-post'] = '/dog';
    attrs['hx-target'] = 'tbody';
    attrs['hx-swap'] = 'afterbegin';
  }

  const selectedDog = dogs.get(selectId);

  return c.html(
      <form hx-disabled-elt="#submit-btn" {...attrs}>
        <div>
          <label for="name">Name</label>
          <input
            id="name"
            name="name"
            required
            size={30}
            type="text"
            value={selectedDog?.name ?? ''} // Access name/breed only if there isn't a selected dog.
            />
        </div>
        <div>
          <label for="breed">Breed</label>
          <input
          id="breed"
          name="breed"
          required
          size={30}
          type="text"
          value={selectedDog?.breed ?? ''}
          />
        </div>
        <div>
          <div class="buttons">
            <button id="submit-btn">{selectId ? 'Update' : 'Add'}</button>
            {selectId && (
                <button hx-put="/deselect" hx-swap="none" type="button">
                  Cancel
                </button>
            )}
          </div>
        </div>
      </form>
  )
})

export default app;

type Dog = {id: string; name: string; breed: string};
const dogs = new Map<string, Dog>();

function addDog(name: string, breed: string): Dog {
  const id = crypto.randomUUID();
  const dog = {id, name, breed};
  dogs.set(id, dog);
  return dog;
}

addDog('Comet', 'Whippet');
addDog('Oscar', 'German Shorthaired Pointer');

let selectId = ''; // id of currently selected dog.

// dogRow tags a dog and returns an HTML table describing it.
function dogRow(dog: Dog, updating = false) {
  // If the dog is being updated, perform an out-of-band swap
  // so the new table row can replace the existing one based on the `id` elt.
  const attrs: {[key: string]: string} = {};
  if (updating) attrs['hx-swap-oob'] = 'true';

  return (
    <tr class="on-hover" id={`row-${dog.id}`} {...attrs}>
      <td>{dog.name}</td>
      <td>{dog.breed}</td>
      <td class="buttons">
        {/*Delete button asks for confirmation first.*/}
        <button
          class="show-on-hover"
          hx-delete={`/dog/${dog.id}`}
          hx-confirm="Are you sure?"
          hx-target="closest tr"
          hx-swap="delete"
          >
          &#x1F5D1;
        </button>
        {/* Select the dog to trigger a selection-change event, causing the form to update
         so the user can modify the name and/or breed of the dog.*/}
        <button
          class="show-on-hover"
          hx-put={'/select/'+dog.id}
          hx-swap="none"
          type="button"
          >
          &#x1F4DD;
        </button>
      </td>
    </tr>
  );
}

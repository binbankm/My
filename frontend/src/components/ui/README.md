# shadcn-ui Components

This directory contains shadcn-ui style components for the ServerPanel frontend.

## Available Components

### Button
A versatile button component with multiple variants and sizes.

**Props:**
- `variant`: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link'
- `size`: 'default' | 'sm' | 'lg' | 'icon'
- `class`: Additional CSS classes

**Usage:**
```vue
<Button variant="default" size="default">Click me</Button>
<Button variant="destructive">Delete</Button>
<Button variant="outline" size="sm">Small Button</Button>
```

### Card Components
A set of components for creating card layouts.

**Card** - The main card container
```vue
<Card>
  <!-- Card content -->
</Card>
```

**CardHeader** - Card header section
```vue
<CardHeader>
  <CardTitle>Title</CardTitle>
  <CardDescription>Description</CardDescription>
</CardHeader>
```

**CardContent** - Main content area
```vue
<CardContent>
  <!-- Your content here -->
</CardContent>
```

**CardFooter** - Footer section
```vue
<CardFooter>
  <Button>Action</Button>
</CardFooter>
```

**Complete Example:**
```vue
<Card>
  <CardHeader>
    <CardTitle>User Profile</CardTitle>
    <CardDescription>Manage your account settings</CardDescription>
  </CardHeader>
  <CardContent>
    <!-- Form fields or content -->
  </CardContent>
  <CardFooter>
    <Button>Save Changes</Button>
  </CardFooter>
</Card>
```

### Input
A styled input field component with support for v-model.

**Props:**
- `type`: Input type (default: 'text')
- `modelValue`: v-model binding
- `class`: Additional CSS classes
- All native HTML input attributes via `v-bind="$attrs"`

**Usage:**
```vue
<Input 
  v-model="username" 
  type="text" 
  placeholder="Enter username"
  required
/>
```

### Label
A label component for form fields.

**Props:**
- `for`: The ID of the associated input
- `class`: Additional CSS classes

**Usage:**
```vue
<Label for="username">Username</Label>
<Input id="username" v-model="username" />
```

## Importing Components

Import components from the index file:

```vue
<script setup>
import { 
  Button, 
  Card, 
  CardHeader, 
  CardTitle, 
  CardContent,
  Input, 
  Label 
} from '@/components/ui'
</script>
```

## Styling

All components use Tailwind CSS with CSS variables defined in `src/assets/index.css`. The components automatically adapt to the color scheme defined in your Tailwind configuration.

## Customization

Each component accepts a `class` prop for additional styling:

```vue
<Button class="mt-4">Custom styled button</Button>
<Card class="border-2 border-blue-500">Custom card</Card>
```

The `cn()` utility function from `@/lib/utils` is used to merge classes properly with Tailwind CSS.

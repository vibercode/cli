# UI Component Examples

## Example Conversations and Expected Responses

### 1. Change Button Color

**User**: "Make the button purple"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "update_theme",
  "data": {
    "colors": {
      "primary": "#8B5CF6"
    }
  },
  "explanation": "Changed primary color to purple for buttons"
}
```

### 2. Add New Input Field

**User**: "Add an input field next to the button"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "add_component",
  "data": {
    "id": "input",
    "type": "atom",
    "name": "EmailInput",
    "category": "Input",
    "properties": {
      "placeholder": "Enter your email",
      "type": "email",
      "label": "Email Address"
    },
    "position": { "x": 540, "y": 193 },
    "size": { "w": 280, "h": 56 }
  },
  "explanation": "Added email input field positioned next to the existing button"
}
```

### 3. Move Component

**User**: "Move the button to the center"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "update_component",
  "data": {
    "id": "button-1752363331264",
    "position": { "x": 400, "y": 300 }
  },
  "explanation": "Moved button to center of the canvas"
}
```

### 4. Change Text Content

**User**: "Change button text to 'Submit'"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "update_component",
  "data": {
    "id": "button-1752363331264",
    "properties": {
      "text": "Submit"
    }
  },
  "explanation": "Updated button text from 'Click me' to 'Submit'"
}
```

### 5. Create Animated Text

**User**: "Add an animated title that says 'Welcome to VibeCode'"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "add_component",
  "data": {
    "id": "animated-text",
    "type": "atom",
    "name": "WelcomeTitle",
    "category": "Typography",
    "properties": {
      "text": "Welcome to VibeCode",
      "effect": "rotate3D",
      "className": "text-white text-4xl font-bold",
      "delay": 0,
      "staggerChildren": true
    },
    "position": { "x": 300, "y": 100 },
    "size": { "w": 400, "h": 48 }
  },
  "explanation": "Added animated welcome title with 3D rotation effect"
}
```

### 6. Create Hero Section

**User**: "Create a hero section with title, subtitle and call-to-action"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "add_component",
  "data": {
    "id": "hero",
    "type": "organism",
    "name": "MainHero",
    "category": "Layout",
    "properties": {
      "title": "Welcome to VibeCode",
      "subtitle": "Build amazing interfaces with AI",
      "ctaText": "Get Started",
      "backgroundImage": "https://images.pexels.com/photos/3861969/pexels-photo-3861969.jpeg?auto=compress&cs=tinysrgb&w=1200"
    },
    "position": { "x": 200, "y": 100 },
    "size": { "w": 900, "h": 480 }
  },
  "explanation": "Created hero section with title, subtitle, and call-to-action button"
}
```

### 7. Create Contact Form

**User**: "Add a contact form with name, email and message fields"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "add_component",
  "data": {
    "id": "form",
    "type": "molecule",
    "name": "ContactForm",
    "category": "Input",
    "properties": {
      "title": "Contact Us",
      "fields": ["name", "email", "message"],
      "submitText": "Send Message"
    },
    "position": { "x": 300, "y": 200 },
    "size": { "w": 400, "h": 360 }
  },
  "explanation": "Added contact form with name, email, and message fields"
}
```

### 8. Add Image Gallery

**User**: "Create a photo gallery with 3 columns"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "add_component",
  "data": {
    "id": "gallery",
    "type": "organism",
    "name": "PhotoGallery",
    "category": "Media",
    "properties": {
      "title": "Photo Gallery",
      "images": [
        "https://images.pexels.com/photos/3861969/pexels-photo-3861969.jpeg?auto=compress&cs=tinysrgb&w=400",
        "https://images.pexels.com/photos/2662116/pexels-photo-2662116.jpeg?auto=compress&cs=tinysrgb&w=400",
        "https://images.pexels.com/photos/1181671/pexels-photo-1181671.jpeg?auto=compress&cs=tinysrgb&w=400"
      ],
      "columns": 3
    },
    "position": { "x": 250, "y": 300 },
    "size": { "w": 640, "h": 400 }
  },
  "explanation": "Created a 3-column photo gallery with sample images"
}
```

### 9. Add Navigation Bar

**User**: "Create a navigation bar with Home, About, Services, Contact"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "add_component",
  "data": {
    "id": "navigation",
    "type": "molecule",
    "name": "MainNavigation",
    "category": "Layout",
    "properties": {
      "items": ["Home", "About", "Services", "Contact"],
      "style": "horizontal"
    },
    "position": { "x": 100, "y": 50 },
    "size": { "w": 800, "h": 60 }
  },
  "explanation": "Added horizontal navigation bar with main menu items"
}
```

### 10. Create Card Component

**User**: "Add a card with image and description"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "add_component",
  "data": {
    "id": "card",
    "type": "molecule",
    "name": "ProductCard",
    "category": "Layout",
    "properties": {
      "title": "Product Card",
      "content": "This is a beautiful product card with image and description",
      "hasImage": true,
      "imageUrl": "https://images.pexels.com/photos/3861969/pexels-photo-3861969.jpeg?auto=compress&cs=tinysrgb&w=400"
    },
    "position": { "x": 300, "y": 200 },
    "size": { "w": 320, "h": 400 }
  },
  "explanation": "Added product card with image, title and description"
}
```

### 11. Add Simple Image

**User**: "Add a rounded profile image"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "add_component",
  "data": {
    "id": "image",
    "type": "atom",
    "name": "ProfileImage",
    "category": "Media",
    "properties": {
      "src": "https://images.pexels.com/photos/3861969/pexels-photo-3861969.jpeg?auto=compress&cs=tinysrgb&w=400",
      "alt": "Profile picture",
      "rounded": true
    },
    "position": { "x": 400, "y": 150 },
    "size": { "w": 200, "h": 200 }
  },
  "explanation": "Added rounded profile image"
}
```

### 12. Update Text Content

**User**: "Change the text to say 'Hello World' and make it larger"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "update_component",
  "data": {
    "id": "text-1234567890",
    "properties": {
      "content": "Hello World",
      "size": "large",
      "weight": "bold"
    }
  },
  "explanation": "Updated text content to 'Hello World' and made it larger and bold"
}
```

### 13. Disable Animations

**User**: "Turn off animations, they're distracting"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "update_theme",
  "data": {
    "effects": {
      "animations": false
    }
  },
  "explanation": "Disabled animations for a cleaner, less distracting interface"
}
```

### 14. Make Everything Bigger

**User**: "Scale up all components"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "update_layout",
  "data": {
    "grid": 8,
    "rowHeight": 80
  },
  "explanation": "Increased grid size and row height to scale up all components"
}
```

### 15. Dark Theme

**User**: "Make it darker"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "update_theme",
  "data": {
    "colors": {
      "background": "#000000",
      "surface": "#111111",
      "text": "#FFFFFF"
    }
  },
  "explanation": "Applied darker background and surface colors for a more dramatic dark theme"
}
```

### 16. Add Translated Text

**User**: "Add a welcome message that can be translated"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "add_component",
  "data": {
    "id": "t-animated",
    "type": "atom",
    "name": "WelcomeMessage",
    "category": "Typography",
    "properties": {
      "id": "welcome.title",
      "className": "text-white text-xl font-medium",
      "effect": "fadeIn",
      "delay": 0,
      "staggerChildren": true
    },
    "position": { "x": 300, "y": 200 },
    "size": { "w": 350, "h": 40 }
  },
  "explanation": "Added translatable welcome message with fade-in animation"
}
```

### 17. Remove Component

**User**: "Delete the button"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "remove_component",
  "data": {
    "id": "button-1752363331264"
  },
  "explanation": "Removed button component from canvas"
}
```

### 18. Remove Multiple Components

**User**: "Remove all text components"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "remove_component",
  "data": {
    "id": "text-1234567890"
  },
  "explanation": "Removed text component (removing one at a time for safety)"
}
```

### 19. Ocean Theme

**User**: "Switch to ocean theme"

**AI Response**:

```json
{
  "type": "ui_update",
  "action": "update_theme",
  "data": {
    "colors": {
      "primary": "#0EA5E9",
      "secondary": "#06B6D4",
      "accent": "#10B981",
      "background": "#0C1A2B",
      "surface": "#1E3A5F",
      "text": "#F0F9FF"
    },
    "effects": {
      "glow": false,
      "gradients": true,
      "animations": true
    }
  },
  "explanation": "Applied ocean theme with blue color palette and gradients"
}
```

{% extends "base.tpl" %}

{% block content %}
{% if error or not recipes %}
<div class="error">{{ errorText or 'Recipe not found.' }}</div>
{% else %}
<div class="permindex">
  {% for recipe in recipes %}
    <div class="recipe">
      <a href="/recipe/{{ recipe.name }}">{{ recipe.title }}</a> 
      <span class="description">{{ recipe.description }}</span>
    </div>
  {% endfor %}
</div>
{% endif %}
{% endblock %}

{% extends "base.tpl" %}

{% block styles %}
{{ super() }}
  <link rel="stylesheet" href="/css/recipe.css">
{% endblock styles %}


{% block content %}
{% if error or not recipe %}
<div class="error">{{ errorText or 'Recipe not found.' }}</div>
{% else %}
<div class="recipe">
  <div class="reference"><span class="name">{{ recipe.name }}</span><span class="category">{{ recipe.category }}</span></div>
  <div class="title">{{ recipe.title }}</div>
  <div class="description">{{ recipe.description }}</div>
  {% if recipe.introduction %}
    <div class="introduction">
      {% for line in recipe.introduction %}
        <div class="line">{{ line|interpret }}</div>
      {% endfor %}
    </div>
  {% endif %}

  {% for section in recipe.sections %}
    <div class="section">
      {% if section.name %}
         <div class="named-section">
	 <span class="name">{{ section.name }}</span>
      {% endif %}

      {% if section.ingredient_sets %}
         <div class="ingredients">
           {% for set in section.ingredient_sets %}
	     <div class="header">INGREDIENTS {% if set.yield %} ({{ chooseUnits(set.yield.us, set.yield.metric) }}){% endif %}</div>
             {% for ig in set.ingredients %}
	       <div class="ingredient">
	       	 <span class="amount">{{ chooseUnits(ig.us, ig.metric) }}</span>
		 <span class="item">{{ig.ingredient}}</span>
	         {% if ig.comments %}
		   {% for comment in ig.comments %}
	             <span class="comment">{{ comment|interpret }}</span>
                   {% endfor %}
		 {% endif %}
	       </div>
             {% endfor %}
           {% endfor %}
	 </div>
      {% endif %}

      {% if section.procedure_sets %}
         <div class="procedures">
           {% for set in section.procedure_sets %}
	     <div class="header">PROCEDURES{% if set.name %} {{ set.name }}{% endif %}</div>
             {% for sk in set.procedures %}
	       <div class="procedure">
		 <ol>
		   <li value="{{ sk.index }}">
		     {% for comment in sk.comments %}
		       <div class="comment">{{ comment|interpret }}</div>
		     {% endfor %}
		   </li>
		 </ol>
	       </div>
             {% endfor %}
           {% endfor %}
	 </div>
      {% endif %}

      {% if section.comments %}
	<div class="comments">
	  {% for comment in section.comments %}
	  <div class="comment">{{ comment|interpret }}</div>
          {% endfor %}
	</div>
      {% endif %}

      {% if section.name %}</div>{% endif %}
    </div>
  {% endfor %}

  {% if recipe.notes %}
    <div class="section">
      <div class="name">NOTES</div>
      <div class="notes">
	{% for line in recipe.notes %}
          <div class="line">{{ line|interpret }}</div>
        {% endfor %}
      </div>
    </div>
  {% endif %}

  {% if recipe.ratings %}
    <div class="section">
      <div class="name">RATINGS</div>
      <div class="notes">
	{% for line in recipe.ratings.comments %}
          <div class="line">{{ line|interpret }}</div>
        {% endfor %}
      </div>
    </div>
  {% endif %}

  {% if recipe.footer %}
    <div class="section">
      <div class="name">CONTRIBUTER</div>
      <div class="contributer">
	{% for line in recipe.footer %}
          <div class="line">{{ line|interpret }}</div>
        {% endfor %}
      </div>
    </div>
  {% endif %}
</div>
{% endif %}
{% endblock %}

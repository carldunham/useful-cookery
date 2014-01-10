{% extends "base.tpl" %}

{% block content %}
{% if error or not recipe %}
<div class="error">{{ errorText or 'Recipe not found.' }}</div>
{% else %}
<div class="recipe">
  <div class=headline><span class="name">{{ recipe.name }}</span></div>
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
      <div class="name">{{ section.name or '(unnamed section)' }}</div>

      {% if section.ingredient_sets %}
         <div class="ingredients">
           {% for set in section.ingredient_sets %}
	     <div class="header">INGREDIENTS: {% if set.yield %} ({{ set.yield.us }}){% endif %}</div>
             {% for ig in set.ingredients %}
	       <div class="ingredient">
		 <div>
		   <span class="amount">{{ig.us}}</span>
		   <span class="item">{{ig.ingredient}}</span>
		 </div>
	         {% if ig.comments %}
		   <div class="comments">
		   {% for comment in ig.comments %}
	             <div class="comment">{{ comment|interpret }}</div>
                   {% endfor %}
		   </div>
		 {% endif %}
	       </div>
             {% endfor %}
           {% endfor %}
	 </div>
      {% endif %}

      {% if section.procedure_sets %}
         <div class="procedures">
           {% for set in section.procedure_sets %}
	     <div class="header">PROCEDURES:</div>
             {% for sk in set.procedures %}
	       <div class="procedure">
		 <div>
		   <span class="index">{{ sk.index }}.</span>
		   <span class="comments">{{ sk.comments[0]|interpret }}</span>
		   {% for comment in sk.comments[1:] %}
		     <span class="index next-line"></span>
		     <span class="comment">{{ comment|interpret }}</span>
		   {% endfor %}
		 </div>
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
      <div class="footer">
	{% for line in recipe.footer %}
          <div class="line">{{ line|interpret }}</div>
        {% endfor %}
      </div>
    </div>
  {% endif %}
</div>
{% endif %}
{% endblock %}
